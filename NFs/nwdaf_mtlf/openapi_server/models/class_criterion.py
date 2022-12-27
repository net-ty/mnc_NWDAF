# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.dispersion_class import DispersionClass
from openapi_server.models.matching_direction import MatchingDirection
from openapi_server import util

from openapi_server.models.dispersion_class import DispersionClass  # noqa: E501
from openapi_server.models.matching_direction import MatchingDirection  # noqa: E501

class ClassCriterion(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, disper_class=None, class_threshold=None, thres_match=None):  # noqa: E501
        """ClassCriterion - a model defined in OpenAPI

        :param disper_class: The disper_class of this ClassCriterion.  # noqa: E501
        :type disper_class: DispersionClass
        :param class_threshold: The class_threshold of this ClassCriterion.  # noqa: E501
        :type class_threshold: int
        :param thres_match: The thres_match of this ClassCriterion.  # noqa: E501
        :type thres_match: MatchingDirection
        """
        self.openapi_types = {
            'disper_class': DispersionClass,
            'class_threshold': int,
            'thres_match': MatchingDirection
        }

        self.attribute_map = {
            'disper_class': 'disperClass',
            'class_threshold': 'classThreshold',
            'thres_match': 'thresMatch'
        }

        self.disper_class = disper_class
        self.class_threshold = class_threshold
        self.thres_match = thres_match

    @classmethod
    def from_dict(cls, dikt) -> 'ClassCriterion':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The ClassCriterion of this ClassCriterion.  # noqa: E501
        :rtype: ClassCriterion
        """
        return util.deserialize_model(dikt, cls)

    @property
    def disper_class(self):
        """Gets the disper_class of this ClassCriterion.


        :return: The disper_class of this ClassCriterion.
        :rtype: DispersionClass
        """
        return self._disper_class

    @disper_class.setter
    def disper_class(self, disper_class):
        """Sets the disper_class of this ClassCriterion.


        :param disper_class: The disper_class of this ClassCriterion.
        :type disper_class: DispersionClass
        """
        if disper_class is None:
            raise ValueError("Invalid value for `disper_class`, must not be `None`")  # noqa: E501

        self._disper_class = disper_class

    @property
    def class_threshold(self):
        """Gets the class_threshold of this ClassCriterion.

        Unsigned integer indicating Sampling Ratio (see clauses 4.15.1 of 3GPP TS 23.502), expressed in percent.    # noqa: E501

        :return: The class_threshold of this ClassCriterion.
        :rtype: int
        """
        return self._class_threshold

    @class_threshold.setter
    def class_threshold(self, class_threshold):
        """Sets the class_threshold of this ClassCriterion.

        Unsigned integer indicating Sampling Ratio (see clauses 4.15.1 of 3GPP TS 23.502), expressed in percent.    # noqa: E501

        :param class_threshold: The class_threshold of this ClassCriterion.
        :type class_threshold: int
        """
        if class_threshold is None:
            raise ValueError("Invalid value for `class_threshold`, must not be `None`")  # noqa: E501
        if class_threshold is not None and class_threshold > 100:  # noqa: E501
            raise ValueError("Invalid value for `class_threshold`, must be a value less than or equal to `100`")  # noqa: E501
        if class_threshold is not None and class_threshold < 1:  # noqa: E501
            raise ValueError("Invalid value for `class_threshold`, must be a value greater than or equal to `1`")  # noqa: E501

        self._class_threshold = class_threshold

    @property
    def thres_match(self):
        """Gets the thres_match of this ClassCriterion.


        :return: The thres_match of this ClassCriterion.
        :rtype: MatchingDirection
        """
        return self._thres_match

    @thres_match.setter
    def thres_match(self, thres_match):
        """Sets the thres_match of this ClassCriterion.


        :param thres_match: The thres_match of this ClassCriterion.
        :type thres_match: MatchingDirection
        """
        if thres_match is None:
            raise ValueError("Invalid value for `thres_match`, must not be `None`")  # noqa: E501

        self._thres_match = thres_match