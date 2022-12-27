# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.ladn_info import LadnInfo
from openapi_server.models.presence_info import PresenceInfo
from openapi_server.models.snssai import Snssai
from openapi_server import util

from openapi_server.models.ladn_info import LadnInfo  # noqa: E501
from openapi_server.models.presence_info import PresenceInfo  # noqa: E501
from openapi_server.models.snssai import Snssai  # noqa: E501

class AmfEventArea(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, presence_info=None, ladn_info=None, s_nssai=None, nsi_id=None):  # noqa: E501
        """AmfEventArea - a model defined in OpenAPI

        :param presence_info: The presence_info of this AmfEventArea.  # noqa: E501
        :type presence_info: PresenceInfo
        :param ladn_info: The ladn_info of this AmfEventArea.  # noqa: E501
        :type ladn_info: LadnInfo
        :param s_nssai: The s_nssai of this AmfEventArea.  # noqa: E501
        :type s_nssai: Snssai
        :param nsi_id: The nsi_id of this AmfEventArea.  # noqa: E501
        :type nsi_id: str
        """
        self.openapi_types = {
            'presence_info': PresenceInfo,
            'ladn_info': LadnInfo,
            's_nssai': Snssai,
            'nsi_id': str
        }

        self.attribute_map = {
            'presence_info': 'presenceInfo',
            'ladn_info': 'ladnInfo',
            's_nssai': 'sNssai',
            'nsi_id': 'nsiId'
        }

        self.presence_info = presence_info
        self.ladn_info = ladn_info
        self.s_nssai = s_nssai
        self.nsi_id = nsi_id

    @classmethod
    def from_dict(cls, dikt) -> 'AmfEventArea':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The AmfEventArea of this AmfEventArea.  # noqa: E501
        :rtype: AmfEventArea
        """
        return util.deserialize_model(dikt, cls)

    @property
    def presence_info(self):
        """Gets the presence_info of this AmfEventArea.


        :return: The presence_info of this AmfEventArea.
        :rtype: PresenceInfo
        """
        return self._presence_info

    @presence_info.setter
    def presence_info(self, presence_info):
        """Sets the presence_info of this AmfEventArea.


        :param presence_info: The presence_info of this AmfEventArea.
        :type presence_info: PresenceInfo
        """

        self._presence_info = presence_info

    @property
    def ladn_info(self):
        """Gets the ladn_info of this AmfEventArea.


        :return: The ladn_info of this AmfEventArea.
        :rtype: LadnInfo
        """
        return self._ladn_info

    @ladn_info.setter
    def ladn_info(self, ladn_info):
        """Sets the ladn_info of this AmfEventArea.


        :param ladn_info: The ladn_info of this AmfEventArea.
        :type ladn_info: LadnInfo
        """

        self._ladn_info = ladn_info

    @property
    def s_nssai(self):
        """Gets the s_nssai of this AmfEventArea.


        :return: The s_nssai of this AmfEventArea.
        :rtype: Snssai
        """
        return self._s_nssai

    @s_nssai.setter
    def s_nssai(self, s_nssai):
        """Sets the s_nssai of this AmfEventArea.


        :param s_nssai: The s_nssai of this AmfEventArea.
        :type s_nssai: Snssai
        """

        self._s_nssai = s_nssai

    @property
    def nsi_id(self):
        """Gets the nsi_id of this AmfEventArea.

        Contains the Identifier of the selected Network Slice instance  # noqa: E501

        :return: The nsi_id of this AmfEventArea.
        :rtype: str
        """
        return self._nsi_id

    @nsi_id.setter
    def nsi_id(self, nsi_id):
        """Sets the nsi_id of this AmfEventArea.

        Contains the Identifier of the selected Network Slice instance  # noqa: E501

        :param nsi_id: The nsi_id of this AmfEventArea.
        :type nsi_id: str
        """

        self._nsi_id = nsi_id
