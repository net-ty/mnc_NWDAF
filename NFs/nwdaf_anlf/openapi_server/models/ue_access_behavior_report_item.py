# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.access_state_transition_type import AccessStateTransitionType
from openapi_server import util

from openapi_server.models.access_state_transition_type import AccessStateTransitionType  # noqa: E501

class UeAccessBehaviorReportItem(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, state_transition_type=None, spacing=None, duration=None):  # noqa: E501
        """UeAccessBehaviorReportItem - a model defined in OpenAPI

        :param state_transition_type: The state_transition_type of this UeAccessBehaviorReportItem.  # noqa: E501
        :type state_transition_type: AccessStateTransitionType
        :param spacing: The spacing of this UeAccessBehaviorReportItem.  # noqa: E501
        :type spacing: int
        :param duration: The duration of this UeAccessBehaviorReportItem.  # noqa: E501
        :type duration: int
        """
        self.openapi_types = {
            'state_transition_type': AccessStateTransitionType,
            'spacing': int,
            'duration': int
        }

        self.attribute_map = {
            'state_transition_type': 'stateTransitionType',
            'spacing': 'spacing',
            'duration': 'duration'
        }

        self.state_transition_type = state_transition_type
        self.spacing = spacing
        self.duration = duration

    @classmethod
    def from_dict(cls, dikt) -> 'UeAccessBehaviorReportItem':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The UeAccessBehaviorReportItem of this UeAccessBehaviorReportItem.  # noqa: E501
        :rtype: UeAccessBehaviorReportItem
        """
        return util.deserialize_model(dikt, cls)

    @property
    def state_transition_type(self):
        """Gets the state_transition_type of this UeAccessBehaviorReportItem.


        :return: The state_transition_type of this UeAccessBehaviorReportItem.
        :rtype: AccessStateTransitionType
        """
        return self._state_transition_type

    @state_transition_type.setter
    def state_transition_type(self, state_transition_type):
        """Sets the state_transition_type of this UeAccessBehaviorReportItem.


        :param state_transition_type: The state_transition_type of this UeAccessBehaviorReportItem.
        :type state_transition_type: AccessStateTransitionType
        """
        if state_transition_type is None:
            raise ValueError("Invalid value for `state_transition_type`, must not be `None`")  # noqa: E501

        self._state_transition_type = state_transition_type

    @property
    def spacing(self):
        """Gets the spacing of this UeAccessBehaviorReportItem.

        indicating a time in seconds.  # noqa: E501

        :return: The spacing of this UeAccessBehaviorReportItem.
        :rtype: int
        """
        return self._spacing

    @spacing.setter
    def spacing(self, spacing):
        """Sets the spacing of this UeAccessBehaviorReportItem.

        indicating a time in seconds.  # noqa: E501

        :param spacing: The spacing of this UeAccessBehaviorReportItem.
        :type spacing: int
        """
        if spacing is None:
            raise ValueError("Invalid value for `spacing`, must not be `None`")  # noqa: E501

        self._spacing = spacing

    @property
    def duration(self):
        """Gets the duration of this UeAccessBehaviorReportItem.

        indicating a time in seconds.  # noqa: E501

        :return: The duration of this UeAccessBehaviorReportItem.
        :rtype: int
        """
        return self._duration

    @duration.setter
    def duration(self, duration):
        """Sets the duration of this UeAccessBehaviorReportItem.

        indicating a time in seconds.  # noqa: E501

        :param duration: The duration of this UeAccessBehaviorReportItem.
        :type duration: int
        """
        if duration is None:
            raise ValueError("Invalid value for `duration`, must not be `None`")  # noqa: E501

        self._duration = duration
