# coding: utf-8

from __future__ import absolute_import
import unittest

from flask import json
from six import BytesIO

from openapi_server.models.nwdaf_ml_model_prov_subsc import NwdafMLModelProvSubsc  # noqa: E501
from openapi_server.models.problem_details import ProblemDetails  # noqa: E501
from openapi_server.test import BaseTestCase


class TestSubscriptionsCollectionController(BaseTestCase):
    """SubscriptionsCollectionController integration test stubs"""

    def test_create_nwdafml_model_provision_subcription(self):
        """Test case for create_nwdafml_model_provision_subcription

        Create a new Individual NWDAF ML Model Provision Subscription resource.
        """
        nwdaf_ml_model_prov_subsc = {
  "notifUri" : "notifUri",
  "notifCorreId" : "notifCorreId",
  "mLEventSubscs" : [ {
    "tgtUe" : {
      "supis" : [ null, null ],
      "gpsis" : [ null, null ],
      "anyUe" : true,
      "intGroupIds" : [ null, null ]
    },
    "mLTargetPeriod" : {
      "startTime" : "2000-01-23T04:56:07.000+00:00",
      "stopTime" : "2000-01-23T04:56:07.000+00:00"
    },
    "mLEventFilter" : {
      "bwRequs" : [ {
        "appId" : "appId",
        "marBwDl" : "marBwDl",
        "mirBwUl" : "mirBwUl",
        "marBwUl" : "marBwUl",
        "mirBwDl" : "mirBwDl"
      }, {
        "appId" : "appId",
        "marBwDl" : "marBwDl",
        "mirBwUl" : "mirBwUl",
        "marBwUl" : "marBwUl",
        "mirBwDl" : "mirBwDl"
      } ],
      "upfId" : "upfId",
      "maxTopAppDlNbr" : 0,
      "excepIds" : [ null, null ],
      "qosRequ" : {
        "gfbrUl" : "gfbrUl",
        "gfbrDl" : "gfbrDl",
        "per" : "per",
        "5qi" : 143,
        "pdb" : 1
      },
      "nwPerfTypes" : [ null, null ],
      "networkArea" : {
        "ncgis" : [ {
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          },
          "nrCellId" : "nrCellId"
        }, {
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          },
          "nrCellId" : "nrCellId"
        } ],
        "tais" : [ {
          "tac" : "tac",
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          }
        }, {
          "tac" : "tac",
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          }
        } ],
        "gRanNodeIds" : [ {
          "eNbId" : "eNbId",
          "wagfId" : "wagfId",
          "tngfId" : "tngfId",
          "gNbId" : {
            "bitLength" : 28,
            "gNBValue" : "gNBValue"
          },
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          },
          "n3IwfId" : "n3IwfId",
          "ngeNbId" : "SMacroNGeNB-34B89"
        }, {
          "eNbId" : "eNbId",
          "wagfId" : "wagfId",
          "tngfId" : "tngfId",
          "gNbId" : {
            "bitLength" : 28,
            "gNBValue" : "gNBValue"
          },
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          },
          "n3IwfId" : "n3IwfId",
          "ngeNbId" : "SMacroNGeNB-34B89"
        } ],
        "ecgis" : [ {
          "eutraCellId" : "eutraCellId",
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          }
        }, {
          "eutraCellId" : "eutraCellId",
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          }
        } ]
      },
      "ratFreqs" : [ {
        "allFreq" : true,
        "freq" : 664727,
        "svcExpThreshold" : {
          "svcExpLevel" : 4.9652185,
          "nfStorageUsage" : 1,
          "avgTrafficRate" : "avgTrafficRate",
          "avgPacketDelay" : 1,
          "congLevel" : 4,
          "nfCpuUsage" : 1,
          "maxTrafficRate" : "maxTrafficRate",
          "maxPacketDelay" : 1,
          "nfMemoryUsage" : 1,
          "nfLoadLevel" : 7,
          "avgPacketLossRate" : 117
        },
        "allRat" : true
      }, {
        "allFreq" : true,
        "freq" : 664727,
        "svcExpThreshold" : {
          "svcExpLevel" : 4.9652185,
          "nfStorageUsage" : 1,
          "avgTrafficRate" : "avgTrafficRate",
          "avgPacketDelay" : 1,
          "congLevel" : 4,
          "nfCpuUsage" : 1,
          "maxTrafficRate" : "maxTrafficRate",
          "maxPacketDelay" : 1,
          "nfMemoryUsage" : 1,
          "nfLoadLevel" : 7,
          "avgPacketLossRate" : 117
        },
        "allRat" : true
      } ],
      "listOfAnaSubsets" : [ null, null ],
      "snssais" : [ {
        "sd" : "sd",
        "sst" : 20
      }, {
        "sd" : "sd",
        "sst" : 20
      } ],
      "dnPerfReqs" : [ {
        "reportThresholds" : [ {
          "svcExpLevel" : 4.9652185,
          "nfStorageUsage" : 1,
          "avgTrafficRate" : "avgTrafficRate",
          "avgPacketDelay" : 1,
          "congLevel" : 4,
          "nfCpuUsage" : 1,
          "maxTrafficRate" : "maxTrafficRate",
          "maxPacketDelay" : 1,
          "nfMemoryUsage" : 1,
          "nfLoadLevel" : 7,
          "avgPacketLossRate" : 117
        }, {
          "svcExpLevel" : 4.9652185,
          "nfStorageUsage" : 1,
          "avgTrafficRate" : "avgTrafficRate",
          "avgPacketDelay" : 1,
          "congLevel" : 4,
          "nfCpuUsage" : 1,
          "maxTrafficRate" : "maxTrafficRate",
          "maxPacketDelay" : 1,
          "nfMemoryUsage" : 1,
          "nfLoadLevel" : 7,
          "avgPacketLossRate" : 117
        } ]
      }, {
        "reportThresholds" : [ {
          "svcExpLevel" : 4.9652185,
          "nfStorageUsage" : 1,
          "avgTrafficRate" : "avgTrafficRate",
          "avgPacketDelay" : 1,
          "congLevel" : 4,
          "nfCpuUsage" : 1,
          "maxTrafficRate" : "maxTrafficRate",
          "maxPacketDelay" : 1,
          "nfMemoryUsage" : 1,
          "nfLoadLevel" : 7,
          "avgPacketLossRate" : 117
        }, {
          "svcExpLevel" : 4.9652185,
          "nfStorageUsage" : 1,
          "avgTrafficRate" : "avgTrafficRate",
          "avgPacketDelay" : 1,
          "congLevel" : 4,
          "nfCpuUsage" : 1,
          "maxTrafficRate" : "maxTrafficRate",
          "maxPacketDelay" : 1,
          "nfMemoryUsage" : 1,
          "nfLoadLevel" : 7,
          "avgPacketLossRate" : 117
        } ]
      } ],
      "dnns" : [ null, null ],
      "ladnDnns" : [ null, null ],
      "redTransReqs" : [ { }, { } ],
      "maxTopAppUlNbr" : 0,
      "appServerAddrs" : [ {
        "fqdn" : "fqdn",
        "ipAddr" : {
          "ipv6Addr" : "2001:db8:85a3::8a2e:370:7334",
          "ipv4Addr" : "198.51.100.1",
          "ipv6Prefix" : "2001:db8:abcd:12::0/64"
        }
      }, {
        "fqdn" : "fqdn",
        "ipAddr" : {
          "ipv6Addr" : "2001:db8:85a3::8a2e:370:7334",
          "ipv4Addr" : "198.51.100.1",
          "ipv6Prefix" : "2001:db8:abcd:12::0/64"
        }
      } ],
      "nfSetIds" : [ null, null ],
      "nsiIdInfos" : [ {
        "snssai" : {
          "sd" : "sd",
          "sst" : 20
        },
        "nsiIds" : [ null, null ]
      }, {
        "snssai" : {
          "sd" : "sd",
          "sst" : 20
        },
        "nsiIds" : [ null, null ]
      } ],
      "disperReqs" : [ {
        "rankCriters" : [ {
          "lowBase" : 93,
          "highBase" : 99
        }, {
          "lowBase" : 93,
          "highBase" : 99
        } ],
        "classCriters" : [ {
          "classThreshold" : 50
        }, {
          "classThreshold" : 50
        } ]
      }, {
        "rankCriters" : [ {
          "lowBase" : 93,
          "highBase" : 99
        }, {
          "lowBase" : 93,
          "highBase" : 99
        } ],
        "classCriters" : [ {
          "classThreshold" : 50
        }, {
          "classThreshold" : 50
        } ]
      } ],
      "exptUeBehav" : {
        "validityTime" : "2000-01-23T04:56:07.000+00:00",
        "communicationDurationTime" : 7,
        "scheduledCommunicationTime" : {
          "timeOfDayEnd" : "timeOfDayEnd",
          "daysOfWeek" : [ null, null, null, null, null ],
          "timeOfDayStart" : "timeOfDayStart"
        },
        "periodicTime" : 9,
        "expectedUmts" : [ {
          "geographicAreas" : [ null, null ],
          "nwAreaInfo" : {
            "ncgis" : [ {
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              },
              "nrCellId" : "nrCellId"
            }, {
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              },
              "nrCellId" : "nrCellId"
            } ],
            "tais" : [ {
              "tac" : "tac",
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              }
            }, {
              "tac" : "tac",
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              }
            } ],
            "gRanNodeIds" : [ {
              "eNbId" : "eNbId",
              "wagfId" : "wagfId",
              "tngfId" : "tngfId",
              "gNbId" : {
                "bitLength" : 28,
                "gNBValue" : "gNBValue"
              },
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              },
              "n3IwfId" : "n3IwfId",
              "ngeNbId" : "SMacroNGeNB-34B89"
            }, {
              "eNbId" : "eNbId",
              "wagfId" : "wagfId",
              "tngfId" : "tngfId",
              "gNbId" : {
                "bitLength" : 28,
                "gNBValue" : "gNBValue"
              },
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              },
              "n3IwfId" : "n3IwfId",
              "ngeNbId" : "SMacroNGeNB-34B89"
            } ],
            "ecgis" : [ {
              "eutraCellId" : "eutraCellId",
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              }
            }, {
              "eutraCellId" : "eutraCellId",
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              }
            } ]
          },
          "civicAddresses" : [ {
            "POBOX" : "POBOX",
            "usageRules" : "usageRules",
            "country" : "country",
            "PRD" : "PRD",
            "PLC" : "PLC",
            "HNO" : "HNO",
            "PRM" : "PRM",
            "HNS" : "HNS",
            "FLR" : "FLR",
            "A1" : "A1",
            "A2" : "A2",
            "A3" : "A3",
            "A4" : "A4",
            "STS" : "STS",
            "A5" : "A5",
            "A6" : "A6",
            "RDSEC" : "RDSEC",
            "providedBy" : "providedBy",
            "LOC" : "LOC",
            "UNIT" : "UNIT",
            "SEAT" : "SEAT",
            "POD" : "POD",
            "RDBR" : "RDBR",
            "method" : "method",
            "LMK" : "LMK",
            "POM" : "POM",
            "ADDCODE" : "ADDCODE",
            "RD" : "RD",
            "PC" : "PC",
            "PCN" : "PCN",
            "NAM" : "NAM",
            "BLD" : "BLD",
            "ROOM" : "ROOM",
            "RDSUBBR" : "RDSUBBR"
          }, {
            "POBOX" : "POBOX",
            "usageRules" : "usageRules",
            "country" : "country",
            "PRD" : "PRD",
            "PLC" : "PLC",
            "HNO" : "HNO",
            "PRM" : "PRM",
            "HNS" : "HNS",
            "FLR" : "FLR",
            "A1" : "A1",
            "A2" : "A2",
            "A3" : "A3",
            "A4" : "A4",
            "STS" : "STS",
            "A5" : "A5",
            "A6" : "A6",
            "RDSEC" : "RDSEC",
            "providedBy" : "providedBy",
            "LOC" : "LOC",
            "UNIT" : "UNIT",
            "SEAT" : "SEAT",
            "POD" : "POD",
            "RDBR" : "RDBR",
            "method" : "method",
            "LMK" : "LMK",
            "POM" : "POM",
            "ADDCODE" : "ADDCODE",
            "RD" : "RD",
            "PC" : "PC",
            "PCN" : "PCN",
            "NAM" : "NAM",
            "BLD" : "BLD",
            "ROOM" : "ROOM",
            "RDSUBBR" : "RDSUBBR"
          } ],
          "umtTime" : {
            "dayOfWeek" : 3,
            "timeOfDay" : "timeOfDay"
          }
        }, {
          "geographicAreas" : [ null, null ],
          "nwAreaInfo" : {
            "ncgis" : [ {
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              },
              "nrCellId" : "nrCellId"
            }, {
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              },
              "nrCellId" : "nrCellId"
            } ],
            "tais" : [ {
              "tac" : "tac",
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              }
            }, {
              "tac" : "tac",
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              }
            } ],
            "gRanNodeIds" : [ {
              "eNbId" : "eNbId",
              "wagfId" : "wagfId",
              "tngfId" : "tngfId",
              "gNbId" : {
                "bitLength" : 28,
                "gNBValue" : "gNBValue"
              },
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              },
              "n3IwfId" : "n3IwfId",
              "ngeNbId" : "SMacroNGeNB-34B89"
            }, {
              "eNbId" : "eNbId",
              "wagfId" : "wagfId",
              "tngfId" : "tngfId",
              "gNbId" : {
                "bitLength" : 28,
                "gNBValue" : "gNBValue"
              },
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              },
              "n3IwfId" : "n3IwfId",
              "ngeNbId" : "SMacroNGeNB-34B89"
            } ],
            "ecgis" : [ {
              "eutraCellId" : "eutraCellId",
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              }
            }, {
              "eutraCellId" : "eutraCellId",
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              }
            } ]
          },
          "civicAddresses" : [ {
            "POBOX" : "POBOX",
            "usageRules" : "usageRules",
            "country" : "country",
            "PRD" : "PRD",
            "PLC" : "PLC",
            "HNO" : "HNO",
            "PRM" : "PRM",
            "HNS" : "HNS",
            "FLR" : "FLR",
            "A1" : "A1",
            "A2" : "A2",
            "A3" : "A3",
            "A4" : "A4",
            "STS" : "STS",
            "A5" : "A5",
            "A6" : "A6",
            "RDSEC" : "RDSEC",
            "providedBy" : "providedBy",
            "LOC" : "LOC",
            "UNIT" : "UNIT",
            "SEAT" : "SEAT",
            "POD" : "POD",
            "RDBR" : "RDBR",
            "method" : "method",
            "LMK" : "LMK",
            "POM" : "POM",
            "ADDCODE" : "ADDCODE",
            "RD" : "RD",
            "PC" : "PC",
            "PCN" : "PCN",
            "NAM" : "NAM",
            "BLD" : "BLD",
            "ROOM" : "ROOM",
            "RDSUBBR" : "RDSUBBR"
          }, {
            "POBOX" : "POBOX",
            "usageRules" : "usageRules",
            "country" : "country",
            "PRD" : "PRD",
            "PLC" : "PLC",
            "HNO" : "HNO",
            "PRM" : "PRM",
            "HNS" : "HNS",
            "FLR" : "FLR",
            "A1" : "A1",
            "A2" : "A2",
            "A3" : "A3",
            "A4" : "A4",
            "STS" : "STS",
            "A5" : "A5",
            "A6" : "A6",
            "RDSEC" : "RDSEC",
            "providedBy" : "providedBy",
            "LOC" : "LOC",
            "UNIT" : "UNIT",
            "SEAT" : "SEAT",
            "POD" : "POD",
            "RDBR" : "RDBR",
            "method" : "method",
            "LMK" : "LMK",
            "POM" : "POM",
            "ADDCODE" : "ADDCODE",
            "RD" : "RD",
            "PC" : "PC",
            "PCN" : "PCN",
            "NAM" : "NAM",
            "BLD" : "BLD",
            "ROOM" : "ROOM",
            "RDSUBBR" : "RDSUBBR"
          } ],
          "umtTime" : {
            "dayOfWeek" : 3,
            "timeOfDay" : "timeOfDay"
          }
        } ],
        "batteryIndication" : {
          "replaceableInd" : true,
          "batteryInd" : true,
          "rechargeableInd" : true
        }
      },
      "wlanReqs" : [ {
        "ssIds" : [ "ssIds", "ssIds" ],
        "bssIds" : [ "bssIds", "bssIds" ]
      }, {
        "ssIds" : [ "ssIds", "ssIds" ],
        "bssIds" : [ "bssIds", "bssIds" ]
      } ],
      "appIds" : [ null, null ],
      "dnais" : [ null, null ],
      "nfTypes" : [ null, null ],
      "visitedAreas" : [ {
        "ncgis" : [ {
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          },
          "nrCellId" : "nrCellId"
        }, {
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          },
          "nrCellId" : "nrCellId"
        } ],
        "tais" : [ {
          "tac" : "tac",
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          }
        }, {
          "tac" : "tac",
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          }
        } ],
        "gRanNodeIds" : [ {
          "eNbId" : "eNbId",
          "wagfId" : "wagfId",
          "tngfId" : "tngfId",
          "gNbId" : {
            "bitLength" : 28,
            "gNBValue" : "gNBValue"
          },
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          },
          "n3IwfId" : "n3IwfId",
          "ngeNbId" : "SMacroNGeNB-34B89"
        }, {
          "eNbId" : "eNbId",
          "wagfId" : "wagfId",
          "tngfId" : "tngfId",
          "gNbId" : {
            "bitLength" : 28,
            "gNBValue" : "gNBValue"
          },
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          },
          "n3IwfId" : "n3IwfId",
          "ngeNbId" : "SMacroNGeNB-34B89"
        } ],
        "ecgis" : [ {
          "eutraCellId" : "eutraCellId",
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          }
        }, {
          "eutraCellId" : "eutraCellId",
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          }
        } ]
      }, {
        "ncgis" : [ {
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          },
          "nrCellId" : "nrCellId"
        }, {
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          },
          "nrCellId" : "nrCellId"
        } ],
        "tais" : [ {
          "tac" : "tac",
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          }
        }, {
          "tac" : "tac",
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          }
        } ],
        "gRanNodeIds" : [ {
          "eNbId" : "eNbId",
          "wagfId" : "wagfId",
          "tngfId" : "tngfId",
          "gNbId" : {
            "bitLength" : 28,
            "gNBValue" : "gNBValue"
          },
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          },
          "n3IwfId" : "n3IwfId",
          "ngeNbId" : "SMacroNGeNB-34B89"
        }, {
          "eNbId" : "eNbId",
          "wagfId" : "wagfId",
          "tngfId" : "tngfId",
          "gNbId" : {
            "bitLength" : 28,
            "gNBValue" : "gNBValue"
          },
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          },
          "n3IwfId" : "n3IwfId",
          "ngeNbId" : "SMacroNGeNB-34B89"
        } ],
        "ecgis" : [ {
          "eutraCellId" : "eutraCellId",
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          }
        }, {
          "eutraCellId" : "eutraCellId",
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          }
        } ]
      } ],
      "anySlice" : true,
      "nfInstanceIds" : [ null, null ]
    },
    "expiryTime" : "2000-01-23T04:56:07.000+00:00"
  }, {
    "tgtUe" : {
      "supis" : [ null, null ],
      "gpsis" : [ null, null ],
      "anyUe" : true,
      "intGroupIds" : [ null, null ]
    },
    "mLTargetPeriod" : {
      "startTime" : "2000-01-23T04:56:07.000+00:00",
      "stopTime" : "2000-01-23T04:56:07.000+00:00"
    },
    "mLEventFilter" : {
      "bwRequs" : [ {
        "appId" : "appId",
        "marBwDl" : "marBwDl",
        "mirBwUl" : "mirBwUl",
        "marBwUl" : "marBwUl",
        "mirBwDl" : "mirBwDl"
      }, {
        "appId" : "appId",
        "marBwDl" : "marBwDl",
        "mirBwUl" : "mirBwUl",
        "marBwUl" : "marBwUl",
        "mirBwDl" : "mirBwDl"
      } ],
      "upfId" : "upfId",
      "maxTopAppDlNbr" : 0,
      "excepIds" : [ null, null ],
      "qosRequ" : {
        "gfbrUl" : "gfbrUl",
        "gfbrDl" : "gfbrDl",
        "per" : "per",
        "5qi" : 143,
        "pdb" : 1
      },
      "nwPerfTypes" : [ null, null ],
      "networkArea" : {
        "ncgis" : [ {
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          },
          "nrCellId" : "nrCellId"
        }, {
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          },
          "nrCellId" : "nrCellId"
        } ],
        "tais" : [ {
          "tac" : "tac",
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          }
        }, {
          "tac" : "tac",
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          }
        } ],
        "gRanNodeIds" : [ {
          "eNbId" : "eNbId",
          "wagfId" : "wagfId",
          "tngfId" : "tngfId",
          "gNbId" : {
            "bitLength" : 28,
            "gNBValue" : "gNBValue"
          },
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          },
          "n3IwfId" : "n3IwfId",
          "ngeNbId" : "SMacroNGeNB-34B89"
        }, {
          "eNbId" : "eNbId",
          "wagfId" : "wagfId",
          "tngfId" : "tngfId",
          "gNbId" : {
            "bitLength" : 28,
            "gNBValue" : "gNBValue"
          },
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          },
          "n3IwfId" : "n3IwfId",
          "ngeNbId" : "SMacroNGeNB-34B89"
        } ],
        "ecgis" : [ {
          "eutraCellId" : "eutraCellId",
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          }
        }, {
          "eutraCellId" : "eutraCellId",
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          }
        } ]
      },
      "ratFreqs" : [ {
        "allFreq" : true,
        "freq" : 664727,
        "svcExpThreshold" : {
          "svcExpLevel" : 4.9652185,
          "nfStorageUsage" : 1,
          "avgTrafficRate" : "avgTrafficRate",
          "avgPacketDelay" : 1,
          "congLevel" : 4,
          "nfCpuUsage" : 1,
          "maxTrafficRate" : "maxTrafficRate",
          "maxPacketDelay" : 1,
          "nfMemoryUsage" : 1,
          "nfLoadLevel" : 7,
          "avgPacketLossRate" : 117
        },
        "allRat" : true
      }, {
        "allFreq" : true,
        "freq" : 664727,
        "svcExpThreshold" : {
          "svcExpLevel" : 4.9652185,
          "nfStorageUsage" : 1,
          "avgTrafficRate" : "avgTrafficRate",
          "avgPacketDelay" : 1,
          "congLevel" : 4,
          "nfCpuUsage" : 1,
          "maxTrafficRate" : "maxTrafficRate",
          "maxPacketDelay" : 1,
          "nfMemoryUsage" : 1,
          "nfLoadLevel" : 7,
          "avgPacketLossRate" : 117
        },
        "allRat" : true
      } ],
      "listOfAnaSubsets" : [ null, null ],
      "snssais" : [ {
        "sd" : "sd",
        "sst" : 20
      }, {
        "sd" : "sd",
        "sst" : 20
      } ],
      "dnPerfReqs" : [ {
        "reportThresholds" : [ {
          "svcExpLevel" : 4.9652185,
          "nfStorageUsage" : 1,
          "avgTrafficRate" : "avgTrafficRate",
          "avgPacketDelay" : 1,
          "congLevel" : 4,
          "nfCpuUsage" : 1,
          "maxTrafficRate" : "maxTrafficRate",
          "maxPacketDelay" : 1,
          "nfMemoryUsage" : 1,
          "nfLoadLevel" : 7,
          "avgPacketLossRate" : 117
        }, {
          "svcExpLevel" : 4.9652185,
          "nfStorageUsage" : 1,
          "avgTrafficRate" : "avgTrafficRate",
          "avgPacketDelay" : 1,
          "congLevel" : 4,
          "nfCpuUsage" : 1,
          "maxTrafficRate" : "maxTrafficRate",
          "maxPacketDelay" : 1,
          "nfMemoryUsage" : 1,
          "nfLoadLevel" : 7,
          "avgPacketLossRate" : 117
        } ]
      }, {
        "reportThresholds" : [ {
          "svcExpLevel" : 4.9652185,
          "nfStorageUsage" : 1,
          "avgTrafficRate" : "avgTrafficRate",
          "avgPacketDelay" : 1,
          "congLevel" : 4,
          "nfCpuUsage" : 1,
          "maxTrafficRate" : "maxTrafficRate",
          "maxPacketDelay" : 1,
          "nfMemoryUsage" : 1,
          "nfLoadLevel" : 7,
          "avgPacketLossRate" : 117
        }, {
          "svcExpLevel" : 4.9652185,
          "nfStorageUsage" : 1,
          "avgTrafficRate" : "avgTrafficRate",
          "avgPacketDelay" : 1,
          "congLevel" : 4,
          "nfCpuUsage" : 1,
          "maxTrafficRate" : "maxTrafficRate",
          "maxPacketDelay" : 1,
          "nfMemoryUsage" : 1,
          "nfLoadLevel" : 7,
          "avgPacketLossRate" : 117
        } ]
      } ],
      "dnns" : [ null, null ],
      "ladnDnns" : [ null, null ],
      "redTransReqs" : [ { }, { } ],
      "maxTopAppUlNbr" : 0,
      "appServerAddrs" : [ {
        "fqdn" : "fqdn",
        "ipAddr" : {
          "ipv6Addr" : "2001:db8:85a3::8a2e:370:7334",
          "ipv4Addr" : "198.51.100.1",
          "ipv6Prefix" : "2001:db8:abcd:12::0/64"
        }
      }, {
        "fqdn" : "fqdn",
        "ipAddr" : {
          "ipv6Addr" : "2001:db8:85a3::8a2e:370:7334",
          "ipv4Addr" : "198.51.100.1",
          "ipv6Prefix" : "2001:db8:abcd:12::0/64"
        }
      } ],
      "nfSetIds" : [ null, null ],
      "nsiIdInfos" : [ {
        "snssai" : {
          "sd" : "sd",
          "sst" : 20
        },
        "nsiIds" : [ null, null ]
      }, {
        "snssai" : {
          "sd" : "sd",
          "sst" : 20
        },
        "nsiIds" : [ null, null ]
      } ],
      "disperReqs" : [ {
        "rankCriters" : [ {
          "lowBase" : 93,
          "highBase" : 99
        }, {
          "lowBase" : 93,
          "highBase" : 99
        } ],
        "classCriters" : [ {
          "classThreshold" : 50
        }, {
          "classThreshold" : 50
        } ]
      }, {
        "rankCriters" : [ {
          "lowBase" : 93,
          "highBase" : 99
        }, {
          "lowBase" : 93,
          "highBase" : 99
        } ],
        "classCriters" : [ {
          "classThreshold" : 50
        }, {
          "classThreshold" : 50
        } ]
      } ],
      "exptUeBehav" : {
        "validityTime" : "2000-01-23T04:56:07.000+00:00",
        "communicationDurationTime" : 7,
        "scheduledCommunicationTime" : {
          "timeOfDayEnd" : "timeOfDayEnd",
          "daysOfWeek" : [ null, null, null, null, null ],
          "timeOfDayStart" : "timeOfDayStart"
        },
        "periodicTime" : 9,
        "expectedUmts" : [ {
          "geographicAreas" : [ null, null ],
          "nwAreaInfo" : {
            "ncgis" : [ {
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              },
              "nrCellId" : "nrCellId"
            }, {
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              },
              "nrCellId" : "nrCellId"
            } ],
            "tais" : [ {
              "tac" : "tac",
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              }
            }, {
              "tac" : "tac",
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              }
            } ],
            "gRanNodeIds" : [ {
              "eNbId" : "eNbId",
              "wagfId" : "wagfId",
              "tngfId" : "tngfId",
              "gNbId" : {
                "bitLength" : 28,
                "gNBValue" : "gNBValue"
              },
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              },
              "n3IwfId" : "n3IwfId",
              "ngeNbId" : "SMacroNGeNB-34B89"
            }, {
              "eNbId" : "eNbId",
              "wagfId" : "wagfId",
              "tngfId" : "tngfId",
              "gNbId" : {
                "bitLength" : 28,
                "gNBValue" : "gNBValue"
              },
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              },
              "n3IwfId" : "n3IwfId",
              "ngeNbId" : "SMacroNGeNB-34B89"
            } ],
            "ecgis" : [ {
              "eutraCellId" : "eutraCellId",
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              }
            }, {
              "eutraCellId" : "eutraCellId",
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              }
            } ]
          },
          "civicAddresses" : [ {
            "POBOX" : "POBOX",
            "usageRules" : "usageRules",
            "country" : "country",
            "PRD" : "PRD",
            "PLC" : "PLC",
            "HNO" : "HNO",
            "PRM" : "PRM",
            "HNS" : "HNS",
            "FLR" : "FLR",
            "A1" : "A1",
            "A2" : "A2",
            "A3" : "A3",
            "A4" : "A4",
            "STS" : "STS",
            "A5" : "A5",
            "A6" : "A6",
            "RDSEC" : "RDSEC",
            "providedBy" : "providedBy",
            "LOC" : "LOC",
            "UNIT" : "UNIT",
            "SEAT" : "SEAT",
            "POD" : "POD",
            "RDBR" : "RDBR",
            "method" : "method",
            "LMK" : "LMK",
            "POM" : "POM",
            "ADDCODE" : "ADDCODE",
            "RD" : "RD",
            "PC" : "PC",
            "PCN" : "PCN",
            "NAM" : "NAM",
            "BLD" : "BLD",
            "ROOM" : "ROOM",
            "RDSUBBR" : "RDSUBBR"
          }, {
            "POBOX" : "POBOX",
            "usageRules" : "usageRules",
            "country" : "country",
            "PRD" : "PRD",
            "PLC" : "PLC",
            "HNO" : "HNO",
            "PRM" : "PRM",
            "HNS" : "HNS",
            "FLR" : "FLR",
            "A1" : "A1",
            "A2" : "A2",
            "A3" : "A3",
            "A4" : "A4",
            "STS" : "STS",
            "A5" : "A5",
            "A6" : "A6",
            "RDSEC" : "RDSEC",
            "providedBy" : "providedBy",
            "LOC" : "LOC",
            "UNIT" : "UNIT",
            "SEAT" : "SEAT",
            "POD" : "POD",
            "RDBR" : "RDBR",
            "method" : "method",
            "LMK" : "LMK",
            "POM" : "POM",
            "ADDCODE" : "ADDCODE",
            "RD" : "RD",
            "PC" : "PC",
            "PCN" : "PCN",
            "NAM" : "NAM",
            "BLD" : "BLD",
            "ROOM" : "ROOM",
            "RDSUBBR" : "RDSUBBR"
          } ],
          "umtTime" : {
            "dayOfWeek" : 3,
            "timeOfDay" : "timeOfDay"
          }
        }, {
          "geographicAreas" : [ null, null ],
          "nwAreaInfo" : {
            "ncgis" : [ {
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              },
              "nrCellId" : "nrCellId"
            }, {
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              },
              "nrCellId" : "nrCellId"
            } ],
            "tais" : [ {
              "tac" : "tac",
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              }
            }, {
              "tac" : "tac",
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              }
            } ],
            "gRanNodeIds" : [ {
              "eNbId" : "eNbId",
              "wagfId" : "wagfId",
              "tngfId" : "tngfId",
              "gNbId" : {
                "bitLength" : 28,
                "gNBValue" : "gNBValue"
              },
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              },
              "n3IwfId" : "n3IwfId",
              "ngeNbId" : "SMacroNGeNB-34B89"
            }, {
              "eNbId" : "eNbId",
              "wagfId" : "wagfId",
              "tngfId" : "tngfId",
              "gNbId" : {
                "bitLength" : 28,
                "gNBValue" : "gNBValue"
              },
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              },
              "n3IwfId" : "n3IwfId",
              "ngeNbId" : "SMacroNGeNB-34B89"
            } ],
            "ecgis" : [ {
              "eutraCellId" : "eutraCellId",
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              }
            }, {
              "eutraCellId" : "eutraCellId",
              "nid" : "nid",
              "plmnId" : {
                "mnc" : "mnc",
                "mcc" : "mcc"
              }
            } ]
          },
          "civicAddresses" : [ {
            "POBOX" : "POBOX",
            "usageRules" : "usageRules",
            "country" : "country",
            "PRD" : "PRD",
            "PLC" : "PLC",
            "HNO" : "HNO",
            "PRM" : "PRM",
            "HNS" : "HNS",
            "FLR" : "FLR",
            "A1" : "A1",
            "A2" : "A2",
            "A3" : "A3",
            "A4" : "A4",
            "STS" : "STS",
            "A5" : "A5",
            "A6" : "A6",
            "RDSEC" : "RDSEC",
            "providedBy" : "providedBy",
            "LOC" : "LOC",
            "UNIT" : "UNIT",
            "SEAT" : "SEAT",
            "POD" : "POD",
            "RDBR" : "RDBR",
            "method" : "method",
            "LMK" : "LMK",
            "POM" : "POM",
            "ADDCODE" : "ADDCODE",
            "RD" : "RD",
            "PC" : "PC",
            "PCN" : "PCN",
            "NAM" : "NAM",
            "BLD" : "BLD",
            "ROOM" : "ROOM",
            "RDSUBBR" : "RDSUBBR"
          }, {
            "POBOX" : "POBOX",
            "usageRules" : "usageRules",
            "country" : "country",
            "PRD" : "PRD",
            "PLC" : "PLC",
            "HNO" : "HNO",
            "PRM" : "PRM",
            "HNS" : "HNS",
            "FLR" : "FLR",
            "A1" : "A1",
            "A2" : "A2",
            "A3" : "A3",
            "A4" : "A4",
            "STS" : "STS",
            "A5" : "A5",
            "A6" : "A6",
            "RDSEC" : "RDSEC",
            "providedBy" : "providedBy",
            "LOC" : "LOC",
            "UNIT" : "UNIT",
            "SEAT" : "SEAT",
            "POD" : "POD",
            "RDBR" : "RDBR",
            "method" : "method",
            "LMK" : "LMK",
            "POM" : "POM",
            "ADDCODE" : "ADDCODE",
            "RD" : "RD",
            "PC" : "PC",
            "PCN" : "PCN",
            "NAM" : "NAM",
            "BLD" : "BLD",
            "ROOM" : "ROOM",
            "RDSUBBR" : "RDSUBBR"
          } ],
          "umtTime" : {
            "dayOfWeek" : 3,
            "timeOfDay" : "timeOfDay"
          }
        } ],
        "batteryIndication" : {
          "replaceableInd" : true,
          "batteryInd" : true,
          "rechargeableInd" : true
        }
      },
      "wlanReqs" : [ {
        "ssIds" : [ "ssIds", "ssIds" ],
        "bssIds" : [ "bssIds", "bssIds" ]
      }, {
        "ssIds" : [ "ssIds", "ssIds" ],
        "bssIds" : [ "bssIds", "bssIds" ]
      } ],
      "appIds" : [ null, null ],
      "dnais" : [ null, null ],
      "nfTypes" : [ null, null ],
      "visitedAreas" : [ {
        "ncgis" : [ {
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          },
          "nrCellId" : "nrCellId"
        }, {
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          },
          "nrCellId" : "nrCellId"
        } ],
        "tais" : [ {
          "tac" : "tac",
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          }
        }, {
          "tac" : "tac",
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          }
        } ],
        "gRanNodeIds" : [ {
          "eNbId" : "eNbId",
          "wagfId" : "wagfId",
          "tngfId" : "tngfId",
          "gNbId" : {
            "bitLength" : 28,
            "gNBValue" : "gNBValue"
          },
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          },
          "n3IwfId" : "n3IwfId",
          "ngeNbId" : "SMacroNGeNB-34B89"
        }, {
          "eNbId" : "eNbId",
          "wagfId" : "wagfId",
          "tngfId" : "tngfId",
          "gNbId" : {
            "bitLength" : 28,
            "gNBValue" : "gNBValue"
          },
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          },
          "n3IwfId" : "n3IwfId",
          "ngeNbId" : "SMacroNGeNB-34B89"
        } ],
        "ecgis" : [ {
          "eutraCellId" : "eutraCellId",
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          }
        }, {
          "eutraCellId" : "eutraCellId",
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          }
        } ]
      }, {
        "ncgis" : [ {
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          },
          "nrCellId" : "nrCellId"
        }, {
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          },
          "nrCellId" : "nrCellId"
        } ],
        "tais" : [ {
          "tac" : "tac",
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          }
        }, {
          "tac" : "tac",
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          }
        } ],
        "gRanNodeIds" : [ {
          "eNbId" : "eNbId",
          "wagfId" : "wagfId",
          "tngfId" : "tngfId",
          "gNbId" : {
            "bitLength" : 28,
            "gNBValue" : "gNBValue"
          },
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          },
          "n3IwfId" : "n3IwfId",
          "ngeNbId" : "SMacroNGeNB-34B89"
        }, {
          "eNbId" : "eNbId",
          "wagfId" : "wagfId",
          "tngfId" : "tngfId",
          "gNbId" : {
            "bitLength" : 28,
            "gNBValue" : "gNBValue"
          },
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          },
          "n3IwfId" : "n3IwfId",
          "ngeNbId" : "SMacroNGeNB-34B89"
        } ],
        "ecgis" : [ {
          "eutraCellId" : "eutraCellId",
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          }
        }, {
          "eutraCellId" : "eutraCellId",
          "nid" : "nid",
          "plmnId" : {
            "mnc" : "mnc",
            "mcc" : "mcc"
          }
        } ]
      } ],
      "anySlice" : true,
      "nfInstanceIds" : [ null, null ]
    },
    "expiryTime" : "2000-01-23T04:56:07.000+00:00"
  } ],
  "suppFeats" : "suppFeats",
  "eventReq" : {
    "partitionCriteria" : [ null, null ],
    "grpRepTime" : 6,
    "monDur" : "2000-01-23T04:56:07.000+00:00",
    "immRep" : true,
    "maxReportNbr" : 0,
    "repPeriod" : 8,
    "sampRatio" : 90
  },
  "mLEventNotifs" : [ {
    "validityPeriod" : {
      "startTime" : "2000-01-23T04:56:07.000+00:00",
      "stopTime" : "2000-01-23T04:56:07.000+00:00"
    },
    "notifCorreId" : "notifCorreId",
    "spatialValidity" : {
      "ncgis" : [ {
        "nid" : "nid",
        "plmnId" : {
          "mnc" : "mnc",
          "mcc" : "mcc"
        },
        "nrCellId" : "nrCellId"
      }, {
        "nid" : "nid",
        "plmnId" : {
          "mnc" : "mnc",
          "mcc" : "mcc"
        },
        "nrCellId" : "nrCellId"
      } ],
      "tais" : [ {
        "tac" : "tac",
        "nid" : "nid",
        "plmnId" : {
          "mnc" : "mnc",
          "mcc" : "mcc"
        }
      }, {
        "tac" : "tac",
        "nid" : "nid",
        "plmnId" : {
          "mnc" : "mnc",
          "mcc" : "mcc"
        }
      } ],
      "gRanNodeIds" : [ {
        "eNbId" : "eNbId",
        "wagfId" : "wagfId",
        "tngfId" : "tngfId",
        "gNbId" : {
          "bitLength" : 28,
          "gNBValue" : "gNBValue"
        },
        "nid" : "nid",
        "plmnId" : {
          "mnc" : "mnc",
          "mcc" : "mcc"
        },
        "n3IwfId" : "n3IwfId",
        "ngeNbId" : "SMacroNGeNB-34B89"
      }, {
        "eNbId" : "eNbId",
        "wagfId" : "wagfId",
        "tngfId" : "tngfId",
        "gNbId" : {
          "bitLength" : 28,
          "gNBValue" : "gNBValue"
        },
        "nid" : "nid",
        "plmnId" : {
          "mnc" : "mnc",
          "mcc" : "mcc"
        },
        "n3IwfId" : "n3IwfId",
        "ngeNbId" : "SMacroNGeNB-34B89"
      } ],
      "ecgis" : [ {
        "eutraCellId" : "eutraCellId",
        "nid" : "nid",
        "plmnId" : {
          "mnc" : "mnc",
          "mcc" : "mcc"
        }
      }, {
        "eutraCellId" : "eutraCellId",
        "nid" : "nid",
        "plmnId" : {
          "mnc" : "mnc",
          "mcc" : "mcc"
        }
      } ]
    },
    "mLFileAddr" : {
      "mlFileFqdn" : "mlFileFqdn",
      "mLModelUrl" : "mLModelUrl"
    }
  }, {
    "validityPeriod" : {
      "startTime" : "2000-01-23T04:56:07.000+00:00",
      "stopTime" : "2000-01-23T04:56:07.000+00:00"
    },
    "notifCorreId" : "notifCorreId",
    "spatialValidity" : {
      "ncgis" : [ {
        "nid" : "nid",
        "plmnId" : {
          "mnc" : "mnc",
          "mcc" : "mcc"
        },
        "nrCellId" : "nrCellId"
      }, {
        "nid" : "nid",
        "plmnId" : {
          "mnc" : "mnc",
          "mcc" : "mcc"
        },
        "nrCellId" : "nrCellId"
      } ],
      "tais" : [ {
        "tac" : "tac",
        "nid" : "nid",
        "plmnId" : {
          "mnc" : "mnc",
          "mcc" : "mcc"
        }
      }, {
        "tac" : "tac",
        "nid" : "nid",
        "plmnId" : {
          "mnc" : "mnc",
          "mcc" : "mcc"
        }
      } ],
      "gRanNodeIds" : [ {
        "eNbId" : "eNbId",
        "wagfId" : "wagfId",
        "tngfId" : "tngfId",
        "gNbId" : {
          "bitLength" : 28,
          "gNBValue" : "gNBValue"
        },
        "nid" : "nid",
        "plmnId" : {
          "mnc" : "mnc",
          "mcc" : "mcc"
        },
        "n3IwfId" : "n3IwfId",
        "ngeNbId" : "SMacroNGeNB-34B89"
      }, {
        "eNbId" : "eNbId",
        "wagfId" : "wagfId",
        "tngfId" : "tngfId",
        "gNbId" : {
          "bitLength" : 28,
          "gNBValue" : "gNBValue"
        },
        "nid" : "nid",
        "plmnId" : {
          "mnc" : "mnc",
          "mcc" : "mcc"
        },
        "n3IwfId" : "n3IwfId",
        "ngeNbId" : "SMacroNGeNB-34B89"
      } ],
      "ecgis" : [ {
        "eutraCellId" : "eutraCellId",
        "nid" : "nid",
        "plmnId" : {
          "mnc" : "mnc",
          "mcc" : "mcc"
        }
      }, {
        "eutraCellId" : "eutraCellId",
        "nid" : "nid",
        "plmnId" : {
          "mnc" : "mnc",
          "mcc" : "mcc"
        }
      } ]
    },
    "mLFileAddr" : {
      "mlFileFqdn" : "mlFileFqdn",
      "mLModelUrl" : "mLModelUrl"
    }
  } ],
  "failEventReports" : [ { }, { } ]
}
        headers = { 
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'Authorization': 'Bearer special-key',
        }
        response = self.client.open(
            '/nnwdaf-mlmodelprovision/v1/subscriptions',
            method='POST',
            headers=headers,
            data=json.dumps(nwdaf_ml_model_prov_subsc),
            content_type='application/json')
        self.assert200(response,
                       'Response body is : ' + response.data.decode('utf-8'))


if __name__ == '__main__':
    unittest.main()
