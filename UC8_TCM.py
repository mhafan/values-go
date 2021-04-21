# -*- coding: utf-8 -*-
"""
Created on Mon Mar  8 14:32:07 2021

@author: celbac
"""

import redis
import pandas as pd
import xlrd
import itertools
import random
from matplotlib import pyplot as plt
import numpy as np

# use excel sheet to read in test cases
# =============================================================================
# Excel file is in the following format
# parameter1    Parameter2      Parameter 3
# Value1 p1     value1 p2       value1 p3
# value2 p1     value2 p2       value2 p3

# Example:
# weight	age    	sex	   drug	        bolus	infusion	unitVd
# 100	    50	    male   Rocuronium	10	    5	        2
# 	        20	    female Casatracurium			
# =============================================================================
excelfile = r'./test_parameters.xlsx' 
xls = pd.read_excel(excelfile)
workbook=xlrd.open_workbook(excelfile)
sheet = workbook.sheet_by_index(0)

# define class instance
# the port number might change, docker container was set up with same port
#    docker run -p 6379:6379 --name nameofcontainer -d redis
TCM = redis.Redis(host='localhost', port=6379, db=0, charset="utf-8", decode_responses=True)
CNT = redis.Redis(host='localhost', port=6379, db=0)
PUMP = redis.Redis(host='localhost', port=6379, db=0)
PATMOD = redis.Redis(host='localhost', port=6379, db=0)
SENSOR = redis.Redis(host='localhost', port=6379, db=0)

# =============================================================================
# Data record for every vm.X.Y has the following attributes:
# weight (number) - weight of the patient in [kg],
# age (number) - age of the patient (years),
# sex - { male, female }
# drug (string) - name of the drug { Rocuronium, Casatracurium, ... }
# bolus (number) [mL]
# infusion (number) [mL/hr]
# unitVd (number) - unit Volume of distribution in [mL/kg]
# absoluteVd (number) - absolute Volume of distribution [mL]
# targetTOF (number),
# targetPTC (number), TOF and PTC target values of regulation,
# EC50 (number) - concentration of the drug casuing 50% effect [ug/mL]
# TOF (number) - output effect in TOF units,
# PTC (number) - similarly,
# mtime (number) - current model time [s]
# cycle (number).
# =============================================================================

# read in parameters and values from excel sheet to create a list of all combinations
entries = [] 
parameters = []
for i in range(sheet.ncols):
    parameters.append(sheet.cell(0,i).value)
    entries.append([])
    for parameter in parameters:
         if sheet.cell(0,i).value == parameter:
             for j in range(1,sheet.nrows):
                  if sheet.cell(j,i).value != "":
                     if type(sheet.cell(j,i).value) is str:
                         entries[i].append(sheet.cell(j,i).value)
                     else:
                         entries[i].append(int(sheet.cell(j,i).value))
                                                 
# creating a list of all possible parameter combinations
combinations = []
combinations = list(itertools.product(*entries))

#%% database
# create all identifiers and set their initial values
data=[]
for i in range(len(combinations)):
    cycle = 1
    identifier = 'vm.tc'+str(i+1)+'.'+str(cycle)
    data.append([])
# set initial values for experiment vm.tc1.1
    for parameter in parameters:
        parameter_idx = parameters.index(parameter)
        TCM.hset(identifier, str(parameter), combinations[i][parameter_idx])
        TCM.hset(identifier, 'cycle', cycle)

# publish in channel vm.tc1.1, this channel needs to be subscribed to by others
# all participants are notified about new simulation experiment vm.tc1.1, which is about to start
    TCM.publish(identifier,'start')


# This is the loop that happens outside of the TCM
    while(cycle<=10):
        # the loop is TCM->CNT->PUMP->PATMOD->SENSOR->TCM.
        # CNT message activates CNT
        TCM.publish(identifier, 'CNT')
        # CNT can not overwrite values set by TCM, it can only set new keys
        CNT.hset(identifier, 'bolus', 50)
        CNT.hset(identifier, 'infusion', 0.5)
        # CNT calls pump
        CNT.publish(identifier, 'PUMP')
        # pump calls patient model
        PUMP.publish(identifier, 'PATMOD')
        # patient model sets some values and calls sensor
        PATMOD.hset(identifier, 'TOF', random.uniform(95,100))
        PATMOD.hset(identifier, 'PTC', 12)
        PATMOD.publish(identifier, 'SENSOR')
        # SENSOR calls TCM
        SENSOR.publish(identifier, 'TCM')
        
        # important for TCM
        # save data
        data[i].append(TCM.hgetall(identifier))
        
        # TCM terminates current simulation experiment or starts the next cycle
        TCM.hset(identifier, 'mtime', 0)
        cycle = cycle+1
        TCM.hset(identifier, 'cycle', cycle)
        
    TCM.publish(identifier,'stop')

# TCM ends here


#%% Analysis
        
# display the saved data of TOF
TOF = []
for i in range(len(data)):
    TOF.append(np.array(float(data[0][i]['TOF'])))
    
    
plt.plot(np.arange(0,len(data)),(TOF))
plt.title('TOF')
plt.show()
    