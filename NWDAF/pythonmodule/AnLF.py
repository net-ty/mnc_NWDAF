from Model import *
import os.path

def AnLF(model,data_num):
    if os.path.isfile('model.h5'):
        model = load_model('model.h5')
    prediction = model.predict(x_test[data_num:data_num+1],batch_size = 1)
    print(np.argmax(prediction))
    return np.argmax(prediction)

