# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server import util


class NfStatus(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, status_registered=None, status_unregistered=None, status_undiscoverable=None):  # noqa: E501
        """NfStatus - a model defined in OpenAPI

        :param status_registered: The status_registered of this NfStatus.  # noqa: E501
        :type status_registered: int
        :param status_unregistered: The status_unregistered of this NfStatus.  # noqa: E501
        :type status_unregistered: int
        :param status_undiscoverable: The status_undiscoverable of this NfStatus.  # noqa: E501
        :type status_undiscoverable: int
        """
        self.openapi_types = {
            'status_registered': int,
            'status_unregistered': int,
            'status_undiscoverable': int
        }

        self.attribute_map = {
            'status_registered': 'statusRegistered',
            'status_unregistered': 'statusUnregistered',
            'status_undiscoverable': 'statusUndiscoverable'
        }

        self.status_registered = status_registered
        self.status_unregistered = status_unregistered
        self.status_undiscoverable = status_undiscoverable

    @classmethod
    def from_dict(cls, dikt) -> 'NfStatus':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The NfStatus of this NfStatus.  # noqa: E501
        :rtype: NfStatus
        """
        return util.deserialize_model(dikt, cls)

    @property
    def status_registered(self):
        """Gets the status_registered of this NfStatus.

        Unsigned integer indicating Sampling Ratio (see clauses 4.15.1 of 3GPP TS 23.502), expressed in percent.    # noqa: E501

        :return: The status_registered of this NfStatus.
        :rtype: int
        """
        return self._status_registered

    @status_registered.setter
    def status_registered(self, status_registered):
        """Sets the status_registered of this NfStatus.

        Unsigned integer indicating Sampling Ratio (see clauses 4.15.1 of 3GPP TS 23.502), expressed in percent.    # noqa: E501

        :param status_registered: The status_registered of this NfStatus.
        :type status_registered: int
        """
        if status_registered is not None and status_registered > 100:  # noqa: E501
            raise ValueError("Invalid value for `status_registered`, must be a value less than or equal to `100`")  # noqa: E501
        if status_registered is not None and status_registered < 1:  # noqa: E501
            raise ValueError("Invalid value for `status_registered`, must be a value greater than or equal to `1`")  # noqa: E501

        self._status_registered = status_registered

    @property
    def status_unregistered(self):
        """Gets the status_unregistered of this NfStatus.

        Unsigned integer indicating Sampling Ratio (see clauses 4.15.1 of 3GPP TS 23.502), expressed in percent.    # noqa: E501

        :return: The status_unregistered of this NfStatus.
        :rtype: int
        """
        return self._status_unregistered

    @status_unregistered.setter
    def status_unregistered(self, status_unregistered):
        """Sets the status_unregistered of this NfStatus.

        Unsigned integer indicating Sampling Ratio (see clauses 4.15.1 of 3GPP TS 23.502), expressed in percent.    # noqa: E501

        :param status_unregistered: The status_unregistered of this NfStatus.
        :type status_unregistered: int
        """
        if status_unregistered is not None and status_unregistered > 100:  # noqa: E501
            raise ValueError("Invalid value for `status_unregistered`, must be a value less than or equal to `100`")  # noqa: E501
        if status_unregistered is not None and status_unregistered < 1:  # noqa: E501
            raise ValueError("Invalid value for `status_unregistered`, must be a value greater than or equal to `1`")  # noqa: E501

        self._status_unregistered = status_unregistered

    @property
    def status_undiscoverable(self):
        """Gets the status_undiscoverable of this NfStatus.

        Unsigned integer indicating Sampling Ratio (see clauses 4.15.1 of 3GPP TS 23.502), expressed in percent.    # noqa: E501

        :return: The status_undiscoverable of this NfStatus.
        :rtype: int
        """
        return self._status_undiscoverable

    @status_undiscoverable.setter
    def status_undiscoverable(self, status_undiscoverable):
        """Sets the status_undiscoverable of this NfStatus.

        Unsigned integer indicating Sampling Ratio (see clauses 4.15.1 of 3GPP TS 23.502), expressed in percent.    # noqa: E501

        :param status_undiscoverable: The status_undiscoverable of this NfStatus.
        :type status_undiscoverable: int
        """
        if status_undiscoverable is not None and status_undiscoverable > 100:  # noqa: E501
            raise ValueError("Invalid value for `status_undiscoverable`, must be a value less than or equal to `100`")  # noqa: E501
        if status_undiscoverable is not None and status_undiscoverable < 1:  # noqa: E501
            raise ValueError("Invalid value for `status_undiscoverable`, must be a value greater than or equal to `1`")  # noqa: E501

        self._status_undiscoverable = status_undiscoverable
