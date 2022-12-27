from tracemalloc import start
import matplotlib.pyplot as plt
import keras
import pandas as pd
import numpy as np
from keras.models import Sequential
from keras.layers import Dense
from keras.layers import LSTM
from keras.layers import Dropout
from keras.layers import *
from sklearn.preprocessing import MinMaxScaler
from sklearn.metrics import mean_squared_error
from sklearn.metrics import mean_absolute_error
from sklearn.model_selection import train_test_split
from keras.callbacks import EarlyStopping
from keras.models import load_model
import time
from sklearn.metrics import mean_squared_error
from sklearn.metrics import accuracy_score


global drop_avg
global drop_var

def inference1():

   df=pd.read_csv("./AnLF/dataset/tun1.csv")
   df.head(5)

   training_set = df.iloc[:1400].values
   test_set = df.iloc[1400:].values

   sc = MinMaxScaler(feature_range = (0, 1))
   training_set_scaled = sc.fit_transform(training_set)

   model = load_model("./AnLF/models/model.h5")

   dataset_train = df.iloc[:1400]
   dataset_test = df.iloc[1400:]
   dataset_total = pd.concat((dataset_train, dataset_test), axis = 0)

   inputs = dataset_test.values
   inputs = inputs.reshape(-1,1)
   inputs = sc.transform(inputs)


   X_test = []

   for i in range(0,40):

      X_test.append(inputs[i:i+40, 0])

   X_test = np.array(X_test)
   X_test = np.reshape(X_test, (X_test.shape[0], X_test.shape[1], 1))

   predicted_drop = model.predict(X_test)
   predicted_drop = sc.inverse_transform(predicted_drop)

   drop_avg = np.mean(predicted_drop)
   drop_var = np.var(predicted_drop)

   real_drop = df.iloc[1600:1640].values

   RMSE = mean_squared_error(real_drop, predicted_drop)**0.5

   avg = str(drop_avg)
   var = str(drop_var) 
   packetdrop = "(" + avg + "," + var +")"

   return packetdrop

#print(accuracy)
   print(RMSE)
   print(drop_avg) 
   print(drop_var)
   print(inf_time)

#plt.plot(np.arange(1400,2000),dataset_test.values, color = "blue", label = "Real drop")
#plt.plot(np.arange(1600,2000),predicted_drop, color = "red", label = "Predicted drop")
#plt.title("Packet drop prediction")
#plt.xlabel("Time")
#plt.ylabel("Packet drop")
#plt.legend()
#plt.show()

# duration = 20, prediction 20s periodically
# To list
#x = predicted_drop.tolist()
#pred = list(itertools.chain(*x))
#drop_avg = sum(pred)/len(pred)