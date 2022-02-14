# Public-5G NWDAF

This code is based on the Free5GC.

The detailed description is under below.


## Pre setting
### NWDAF Structure
NWDAF (Network Data Analytics Function) is consist of two part: 1) go module; 2) python module.

Go module can be run on "nwdaf.go" which located in "nwdaf" folder.

Python module can be run on "main.py" which located in "nwdaf/pythonmodule" folder.

In addition, before run this NWDAF function,

#### "nwdaf.cfg" should be located in "free5GC/config" folder.

Also the python module is required to 
```
PYTHON = 3.7
```
and 
```
Tensorflow >= 2.0
```
### temp_requester Structure
temp_requester is the requester function which can be using on other NFs. 

If you want to call NWDAF from other NFs, the function in this requester can be used.


## RUN NWDAF
After Located in the function, the running this program is as follows.


##### 1) run NRF Function (which is in Free5GC)


##### 2) run NWDAF go module
```
cd /move/to/your/path/nwdaf
go run nwdaf.go
```

##### 3) run NWDAF python module
```
cd /move/to/your/path/nwdaf/pythonmodule
python main.py 
```

##### 4) run temp_requester
```
cd /move/to/your/path/nwdaf/temp_requester
go run temp requester
```
After that, you should select your number
If "1" is selected, MTLF (model training function) is executed.
Otherwise, "2" is selected, AnLF (analytics function) is executed.
Then, you can try to select a number which means the dataset number.
Now, we using the EMNIST dataset, which is in the python module.
In temp_requester, the image is not transmitted (using the json, the data number is transmitted).



