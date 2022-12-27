# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server import util


class TwapId(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, ss_id=None, bss_id=None, civic_address=None):  # noqa: E501
        """TwapId - a model defined in OpenAPI

        :param ss_id: The ss_id of this TwapId.  # noqa: E501
        :type ss_id: str
        :param bss_id: The bss_id of this TwapId.  # noqa: E501
        :type bss_id: str
        :param civic_address: The civic_address of this TwapId.  # noqa: E501
        :type civic_address: str
        """
        self.openapi_types = {
            'ss_id': str,
            'bss_id': str,
            'civic_address': str
        }

        self.attribute_map = {
            'ss_id': 'ssId',
            'bss_id': 'bssId',
            'civic_address': 'civicAddress'
        }

        self.ss_id = ss_id
        self.bss_id = bss_id
        self.civic_address = civic_address

    @classmethod
    def from_dict(cls, dikt) -> 'TwapId':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The TwapId of this TwapId.  # noqa: E501
        :rtype: TwapId
        """
        return util.deserialize_model(dikt, cls)

    @property
    def ss_id(self):
        """Gets the ss_id of this TwapId.

        This IE shall contain the SSID of the access point to which the UE is attached, that is received over NGAP, see IEEE Std 802.11-2012.    # noqa: E501

        :return: The ss_id of this TwapId.
        :rtype: str
        """
        return self._ss_id

    @ss_id.setter
    def ss_id(self, ss_id):
        """Sets the ss_id of this TwapId.

        This IE shall contain the SSID of the access point to which the UE is attached, that is received over NGAP, see IEEE Std 802.11-2012.    # noqa: E501

        :param ss_id: The ss_id of this TwapId.
        :type ss_id: str
        """
        if ss_id is None:
            raise ValueError("Invalid value for `ss_id`, must not be `None`")  # noqa: E501

        self._ss_id = ss_id

    @property
    def bss_id(self):
        """Gets the bss_id of this TwapId.

        When present, it shall contain the BSSID of the access point to which the UE is attached, for trusted WLAN access, see IEEE Std 802.11-2012.    # noqa: E501

        :return: The bss_id of this TwapId.
        :rtype: str
        """
        return self._bss_id

    @bss_id.setter
    def bss_id(self, bss_id):
        """Sets the bss_id of this TwapId.

        When present, it shall contain the BSSID of the access point to which the UE is attached, for trusted WLAN access, see IEEE Std 802.11-2012.    # noqa: E501

        :param bss_id: The bss_id of this TwapId.
        :type bss_id: str
        """

        self._bss_id = bss_id

    @property
    def civic_address(self):
        """Gets the civic_address of this TwapId.

        string with format 'bytes' as defined in OpenAPI  # noqa: E501

        :return: The civic_address of this TwapId.
        :rtype: str
        """
        return self._civic_address

    @civic_address.setter
    def civic_address(self, civic_address):
        """Sets the civic_address of this TwapId.

        string with format 'bytes' as defined in OpenAPI  # noqa: E501

        :param civic_address: The civic_address of this TwapId.
        :type civic_address: str
        """

        self._civic_address = civic_address
