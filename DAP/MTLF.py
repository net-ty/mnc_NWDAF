from Model import *
import os.path

def MTLF(model):
    if os.path.isfile('model.h5'):
        model = load_model('model.h5')
    model.fit(x_train, y_train, epochs=1)
    model.save('model.h5')
    print("trainig finish")
