# mnc_Private-5G
테스트베드에서 구현할 구조 설명 및 free5gc와 연동방법 설명
![figure 1](https://user-images.githubusercontent.com/88416778/130419189-8b2debbb-1090-45da-bf32-c29ebb3b4f2b.png)


//////////////////////////////////////////////////////////////////////////////////////////

NWDAF (Network Data Analytics Function) analyzes data through intelligent technologies and provides analysis results to other 5G core network functions

Running nefserver and Pythonmodule should be preceded to run modeltraining.go / model_inference.go

The server automatically turns off after timeouts. That time can be changed as we want

//////////////////////////////////////////////////////////////////////////////////////////

NEF (Network exposure function) stores received information as structured data and exposes it to other network functions

This nefserver.go receives and transfers data from/to NWDAF and DAP

The server should be running continuously to deliver data

//////////////////////////////////////////////////////////////////////////////////////////

DAP (Data Analytic Platform) has a python module to train the data from NWDAF and provides meaningful results to NWDAF through NEF




