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

df=pd.read_csv("tun1.csv")
df.head(5)

training_set = df.iloc[:1400].values
test_set = df.iloc[1400:].values
print(training_set.shape)
print(test_set.shape)
sc = MinMaxScaler(feature_range = (0, 1))
training_set_scaled = sc.fit_transform(training_set)

############################################################

X_train = []
y_train = []

for i in range(200, 1400):

   X_train.append(training_set_scaled[i-200:i, 0])
   y_train.append(training_set_scaled[i, 0])

X_train, y_train = np.array(X_train), np.array(y_train)
X_train = np.reshape(X_train, (X_train.shape[0], X_train.shape[1], 1))

model = Sequential()

model.add(LSTM(units = 50, return_sequences = True, input_shape = (X_train.shape[1], 1)))
model.add(Dropout(0.2))
model.add(LSTM(units = 50, return_sequences = True))
model.add(Dropout(0.2))
model.add(LSTM(units = 50, return_sequences = True))
model.add(Dropout(0.2))
model.add(LSTM(units = 50))
model.add(Dropout(0.2))

model.add(Dense(units = 1))
model.compile(optimizer = 'adam', loss = 'mean_squared_error')

model.fit(X_train, y_train, epochs = 200, batch_size = 32)

model.save("model1.h5")

dataset_train = df.iloc[:1400]
dataset_test = df.iloc[1400:]
dataset_total = pd.concat((dataset_train, dataset_test), axis = 0)

inputs = dataset_test.values
inputs = inputs.reshape(-1,1)
inputs = sc.transform(inputs)

X_test = []

for i in range(200, 600):

   X_test.append(inputs[i-200:i, 0])

X_test = np.array(X_test)
X_test = np.reshape(X_test, (X_test.shape[0], X_test.shape[1], 1))

predicted_delay = model.predict(X_test)
predicted_delay = sc.inverse_transform(predicted_delay)

#plt.plot(df.loc[1400:, "delay"],dataset_test.values, color = "red", label = "Real delay")
#plt.plot(df.loc[1600:, "delay"],predicted_delay, color = "blue", label = "Predicted delay")
#plt.xticks(np.arange(0,2000,200))

plt.plot(np.arange(1400,2000),dataset_test.values, color = "blue", label = "Real delay")
plt.plot(np.arange(1600,2000),predicted_delay, color = "red", label = "Predicted delay")
plt.title("Packet delay prediction")
plt.xlabel("Time")
plt.ylabel("Packet delay")
plt.legend()
plt.show()