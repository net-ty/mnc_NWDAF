# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.access_type import AccessType
from openapi_server.models.cm_state import CmState
from openapi_server import util

from openapi_server.models.access_type import AccessType  # noqa: E501
from openapi_server.models.cm_state import CmState  # noqa: E501

class CmInfo(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, cm_state=None, access_type=None):  # noqa: E501
        """CmInfo - a model defined in OpenAPI

        :param cm_state: The cm_state of this CmInfo.  # noqa: E501
        :type cm_state: CmState
        :param access_type: The access_type of this CmInfo.  # noqa: E501
        :type access_type: AccessType
        """
        self.openapi_types = {
            'cm_state': CmState,
            'access_type': AccessType
        }

        self.attribute_map = {
            'cm_state': 'cmState',
            'access_type': 'accessType'
        }

        self.cm_state = cm_state
        self.access_type = access_type

    @classmethod
    def from_dict(cls, dikt) -> 'CmInfo':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The CmInfo of this CmInfo.  # noqa: E501
        :rtype: CmInfo
        """
        return util.deserialize_model(dikt, cls)

    @property
    def cm_state(self):
        """Gets the cm_state of this CmInfo.


        :return: The cm_state of this CmInfo.
        :rtype: CmState
        """
        return self._cm_state

    @cm_state.setter
    def cm_state(self, cm_state):
        """Sets the cm_state of this CmInfo.


        :param cm_state: The cm_state of this CmInfo.
        :type cm_state: CmState
        """
        if cm_state is None:
            raise ValueError("Invalid value for `cm_state`, must not be `None`")  # noqa: E501

        self._cm_state = cm_state

    @property
    def access_type(self):
        """Gets the access_type of this CmInfo.


        :return: The access_type of this CmInfo.
        :rtype: AccessType
        """
        return self._access_type

    @access_type.setter
    def access_type(self, access_type):
        """Sets the access_type of this CmInfo.


        :param access_type: The access_type of this CmInfo.
        :type access_type: AccessType
        """
        if access_type is None:
            raise ValueError("Invalid value for `access_type`, must not be `None`")  # noqa: E501

        self._access_type = access_type
