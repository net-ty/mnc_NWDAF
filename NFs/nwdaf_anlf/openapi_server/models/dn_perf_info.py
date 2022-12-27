# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.dn_perf import DnPerf
from openapi_server.models.snssai import Snssai
from openapi_server import util

from openapi_server.models.dn_perf import DnPerf  # noqa: E501
from openapi_server.models.snssai import Snssai  # noqa: E501

class DnPerfInfo(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, app_id=None, dnn=None, snssai=None, dn_perf=None, confidence=None):  # noqa: E501
        """DnPerfInfo - a model defined in OpenAPI

        :param app_id: The app_id of this DnPerfInfo.  # noqa: E501
        :type app_id: str
        :param dnn: The dnn of this DnPerfInfo.  # noqa: E501
        :type dnn: str
        :param snssai: The snssai of this DnPerfInfo.  # noqa: E501
        :type snssai: Snssai
        :param dn_perf: The dn_perf of this DnPerfInfo.  # noqa: E501
        :type dn_perf: List[DnPerf]
        :param confidence: The confidence of this DnPerfInfo.  # noqa: E501
        :type confidence: int
        """
        self.openapi_types = {
            'app_id': str,
            'dnn': str,
            'snssai': Snssai,
            'dn_perf': List[DnPerf],
            'confidence': int
        }

        self.attribute_map = {
            'app_id': 'appId',
            'dnn': 'dnn',
            'snssai': 'snssai',
            'dn_perf': 'dnPerf',
            'confidence': 'confidence'
        }

        self.app_id = app_id
        self.dnn = dnn
        self.snssai = snssai
        self.dn_perf = dn_perf
        self.confidence = confidence

    @classmethod
    def from_dict(cls, dikt) -> 'DnPerfInfo':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The DnPerfInfo of this DnPerfInfo.  # noqa: E501
        :rtype: DnPerfInfo
        """
        return util.deserialize_model(dikt, cls)

    @property
    def app_id(self):
        """Gets the app_id of this DnPerfInfo.

        String providing an application identifier.  # noqa: E501

        :return: The app_id of this DnPerfInfo.
        :rtype: str
        """
        return self._app_id

    @app_id.setter
    def app_id(self, app_id):
        """Sets the app_id of this DnPerfInfo.

        String providing an application identifier.  # noqa: E501

        :param app_id: The app_id of this DnPerfInfo.
        :type app_id: str
        """

        self._app_id = app_id

    @property
    def dnn(self):
        """Gets the dnn of this DnPerfInfo.

        String representing a Data Network as defined in clause 9A of 3GPP TS 23.003;  it shall contain either a DNN Network Identifier, or a full DNN with both the Network  Identifier and Operator Identifier, as specified in 3GPP TS 23.003 clause 9.1.1 and 9.1.2. It shall be coded as string in which the labels are separated by dots  (e.g. \"Label1.Label2.Label3\").   # noqa: E501

        :return: The dnn of this DnPerfInfo.
        :rtype: str
        """
        return self._dnn

    @dnn.setter
    def dnn(self, dnn):
        """Sets the dnn of this DnPerfInfo.

        String representing a Data Network as defined in clause 9A of 3GPP TS 23.003;  it shall contain either a DNN Network Identifier, or a full DNN with both the Network  Identifier and Operator Identifier, as specified in 3GPP TS 23.003 clause 9.1.1 and 9.1.2. It shall be coded as string in which the labels are separated by dots  (e.g. \"Label1.Label2.Label3\").   # noqa: E501

        :param dnn: The dnn of this DnPerfInfo.
        :type dnn: str
        """

        self._dnn = dnn

    @property
    def snssai(self):
        """Gets the snssai of this DnPerfInfo.


        :return: The snssai of this DnPerfInfo.
        :rtype: Snssai
        """
        return self._snssai

    @snssai.setter
    def snssai(self, snssai):
        """Sets the snssai of this DnPerfInfo.


        :param snssai: The snssai of this DnPerfInfo.
        :type snssai: Snssai
        """

        self._snssai = snssai

    @property
    def dn_perf(self):
        """Gets the dn_perf of this DnPerfInfo.


        :return: The dn_perf of this DnPerfInfo.
        :rtype: List[DnPerf]
        """
        return self._dn_perf

    @dn_perf.setter
    def dn_perf(self, dn_perf):
        """Sets the dn_perf of this DnPerfInfo.


        :param dn_perf: The dn_perf of this DnPerfInfo.
        :type dn_perf: List[DnPerf]
        """
        if dn_perf is None:
            raise ValueError("Invalid value for `dn_perf`, must not be `None`")  # noqa: E501
        if dn_perf is not None and len(dn_perf) < 1:
            raise ValueError("Invalid value for `dn_perf`, number of items must be greater than or equal to `1`")  # noqa: E501

        self._dn_perf = dn_perf

    @property
    def confidence(self):
        """Gets the confidence of this DnPerfInfo.

        Unsigned Integer, i.e. only value 0 and integers above 0 are permissible.  # noqa: E501

        :return: The confidence of this DnPerfInfo.
        :rtype: int
        """
        return self._confidence

    @confidence.setter
    def confidence(self, confidence):
        """Sets the confidence of this DnPerfInfo.

        Unsigned Integer, i.e. only value 0 and integers above 0 are permissible.  # noqa: E501

        :param confidence: The confidence of this DnPerfInfo.
        :type confidence: int
        """
        if confidence is not None and confidence < 0:  # noqa: E501
            raise ValueError("Invalid value for `confidence`, must be a value greater than or equal to `0`")  # noqa: E501

        self._confidence = confidence
