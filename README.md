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




Note: The code available in this repo implements a NWDAF and associated testbed that is presented in the paper:

T. Kim et al., "An Implementation Study of Network Data Analytic Function in 5G," 2022 IEEE International Conference on Consumer Electronics (ICCE), 2022, pp. 1-3, doi: 10.1109/ICCE53296.2022.9730290.Abstract: Network automation and intelligence are evolutionary directions in 5G, and network data analytic function (NWDAF) plays a key role to realize this vision. In this work, we present an implementation result of NWDAF in free5GC that is an open software for 3GPP mobile core networks. The implemented NWDAF module consists of 1) model training logical function (MTLF) to train the model and 2) analytics logic function (AnLF) to provide analytic results based on the trained model. We have verified the operability of NWDAF and released it through Github. Extensive experimental study will be conducted in our future work.[URL](https://ieeexplore.ieee.org/stamp/stamp.jsp?tp=&arnumber=9730290&isnumber=9730121)
