# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.geographical_coordinates import GeographicalCoordinates
from openapi_server import util

from openapi_server.models.geographical_coordinates import GeographicalCoordinates  # noqa: E501

class PointUncertaintyCircleAllOf(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, point=None, uncertainty=None):  # noqa: E501
        """PointUncertaintyCircleAllOf - a model defined in OpenAPI

        :param point: The point of this PointUncertaintyCircleAllOf.  # noqa: E501
        :type point: GeographicalCoordinates
        :param uncertainty: The uncertainty of this PointUncertaintyCircleAllOf.  # noqa: E501
        :type uncertainty: float
        """
        self.openapi_types = {
            'point': GeographicalCoordinates,
            'uncertainty': float
        }

        self.attribute_map = {
            'point': 'point',
            'uncertainty': 'uncertainty'
        }

        self.point = point
        self.uncertainty = uncertainty

    @classmethod
    def from_dict(cls, dikt) -> 'PointUncertaintyCircleAllOf':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The PointUncertaintyCircle_allOf of this PointUncertaintyCircleAllOf.  # noqa: E501
        :rtype: PointUncertaintyCircleAllOf
        """
        return util.deserialize_model(dikt, cls)

    @property
    def point(self):
        """Gets the point of this PointUncertaintyCircleAllOf.


        :return: The point of this PointUncertaintyCircleAllOf.
        :rtype: GeographicalCoordinates
        """
        return self._point

    @point.setter
    def point(self, point):
        """Sets the point of this PointUncertaintyCircleAllOf.


        :param point: The point of this PointUncertaintyCircleAllOf.
        :type point: GeographicalCoordinates
        """
        if point is None:
            raise ValueError("Invalid value for `point`, must not be `None`")  # noqa: E501

        self._point = point

    @property
    def uncertainty(self):
        """Gets the uncertainty of this PointUncertaintyCircleAllOf.

        Indicates value of uncertainty.  # noqa: E501

        :return: The uncertainty of this PointUncertaintyCircleAllOf.
        :rtype: float
        """
        return self._uncertainty

    @uncertainty.setter
    def uncertainty(self, uncertainty):
        """Sets the uncertainty of this PointUncertaintyCircleAllOf.

        Indicates value of uncertainty.  # noqa: E501

        :param uncertainty: The uncertainty of this PointUncertaintyCircleAllOf.
        :type uncertainty: float
        """
        if uncertainty is None:
            raise ValueError("Invalid value for `uncertainty`, must not be `None`")  # noqa: E501
        if uncertainty is not None and uncertainty < 0:  # noqa: E501
            raise ValueError("Invalid value for `uncertainty`, must be a value greater than or equal to `0`")  # noqa: E501

        self._uncertainty = uncertainty
