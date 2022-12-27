# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.ue_trajectory_collection import UeTrajectoryCollection
import re
from openapi_server import util

from openapi_server.models.ue_trajectory_collection import UeTrajectoryCollection  # noqa: E501
import re  # noqa: E501

class UeMobilityCollection(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, gpsi=None, supi=None, app_id=None, ue_trajs=None):  # noqa: E501
        """UeMobilityCollection - a model defined in OpenAPI

        :param gpsi: The gpsi of this UeMobilityCollection.  # noqa: E501
        :type gpsi: str
        :param supi: The supi of this UeMobilityCollection.  # noqa: E501
        :type supi: str
        :param app_id: The app_id of this UeMobilityCollection.  # noqa: E501
        :type app_id: str
        :param ue_trajs: The ue_trajs of this UeMobilityCollection.  # noqa: E501
        :type ue_trajs: List[UeTrajectoryCollection]
        """
        self.openapi_types = {
            'gpsi': str,
            'supi': str,
            'app_id': str,
            'ue_trajs': List[UeTrajectoryCollection]
        }

        self.attribute_map = {
            'gpsi': 'gpsi',
            'supi': 'supi',
            'app_id': 'appId',
            'ue_trajs': 'ueTrajs'
        }

        self.gpsi = gpsi
        self.supi = supi
        self.app_id = app_id
        self.ue_trajs = ue_trajs

    @classmethod
    def from_dict(cls, dikt) -> 'UeMobilityCollection':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The UeMobilityCollection of this UeMobilityCollection.  # noqa: E501
        :rtype: UeMobilityCollection
        """
        return util.deserialize_model(dikt, cls)

    @property
    def gpsi(self):
        """Gets the gpsi of this UeMobilityCollection.

        String identifying a Gpsi shall contain either an External Id or an MSISDN.  It shall be formatted as follows -External Identifier= \"extid-'extid', where 'extid'  shall be formatted according to clause 19.7.2 of 3GPP TS 23.003 that describes an  External Identifier.    # noqa: E501

        :return: The gpsi of this UeMobilityCollection.
        :rtype: str
        """
        return self._gpsi

    @gpsi.setter
    def gpsi(self, gpsi):
        """Sets the gpsi of this UeMobilityCollection.

        String identifying a Gpsi shall contain either an External Id or an MSISDN.  It shall be formatted as follows -External Identifier= \"extid-'extid', where 'extid'  shall be formatted according to clause 19.7.2 of 3GPP TS 23.003 that describes an  External Identifier.    # noqa: E501

        :param gpsi: The gpsi of this UeMobilityCollection.
        :type gpsi: str
        """
        if gpsi is not None and not re.search(r'^(msisdn-[0-9]{5,15}|extid-[^@]+@[^@]+|.+)$', gpsi):  # noqa: E501
            raise ValueError("Invalid value for `gpsi`, must be a follow pattern or equal to `/^(msisdn-[0-9]{5,15}|extid-[^@]+@[^@]+|.+)$/`")  # noqa: E501

        self._gpsi = gpsi

    @property
    def supi(self):
        """Gets the supi of this UeMobilityCollection.

        String identifying a Supi that shall contain either an IMSI, a network specific identifier, a Global Cable Identifier (GCI) or a Global Line Identifier (GLI) as specified in clause  2.2A of 3GPP TS 23.003. It shall be formatted as follows  - for an IMSI \"imsi-<imsi>\", where <imsi> shall be formatted according to clause 2.2    of 3GPP TS 23.003 that describes an IMSI.  - for a network specific identifier \"nai-<nai>, where <nai> shall be formatted    according to clause 28.7.2 of 3GPP TS 23.003 that describes an NAI.  - for a GCI \"gci-<gci>\", where <gci> shall be formatted according to clause 28.15.2    of 3GPP TS 23.003.  - for a GLI \"gli-<gli>\", where <gli> shall be formatted according to clause 28.16.2 of    3GPP TS 23.003.To enable that the value is used as part of an URI, the string shall    only contain characters allowed according to the \"lower-with-hyphen\" naming convention    defined in 3GPP TS 29.501.   # noqa: E501

        :return: The supi of this UeMobilityCollection.
        :rtype: str
        """
        return self._supi

    @supi.setter
    def supi(self, supi):
        """Sets the supi of this UeMobilityCollection.

        String identifying a Supi that shall contain either an IMSI, a network specific identifier, a Global Cable Identifier (GCI) or a Global Line Identifier (GLI) as specified in clause  2.2A of 3GPP TS 23.003. It shall be formatted as follows  - for an IMSI \"imsi-<imsi>\", where <imsi> shall be formatted according to clause 2.2    of 3GPP TS 23.003 that describes an IMSI.  - for a network specific identifier \"nai-<nai>, where <nai> shall be formatted    according to clause 28.7.2 of 3GPP TS 23.003 that describes an NAI.  - for a GCI \"gci-<gci>\", where <gci> shall be formatted according to clause 28.15.2    of 3GPP TS 23.003.  - for a GLI \"gli-<gli>\", where <gli> shall be formatted according to clause 28.16.2 of    3GPP TS 23.003.To enable that the value is used as part of an URI, the string shall    only contain characters allowed according to the \"lower-with-hyphen\" naming convention    defined in 3GPP TS 29.501.   # noqa: E501

        :param supi: The supi of this UeMobilityCollection.
        :type supi: str
        """
        if supi is not None and not re.search(r'^(imsi-[0-9]{5,15}|nai-.+|gci-.+|gli-.+|.+)$', supi):  # noqa: E501
            raise ValueError("Invalid value for `supi`, must be a follow pattern or equal to `/^(imsi-[0-9]{5,15}|nai-.+|gci-.+|gli-.+|.+)$/`")  # noqa: E501

        self._supi = supi

    @property
    def app_id(self):
        """Gets the app_id of this UeMobilityCollection.

        String providing an application identifier.  # noqa: E501

        :return: The app_id of this UeMobilityCollection.
        :rtype: str
        """
        return self._app_id

    @app_id.setter
    def app_id(self, app_id):
        """Sets the app_id of this UeMobilityCollection.

        String providing an application identifier.  # noqa: E501

        :param app_id: The app_id of this UeMobilityCollection.
        :type app_id: str
        """
        if app_id is None:
            raise ValueError("Invalid value for `app_id`, must not be `None`")  # noqa: E501

        self._app_id = app_id

    @property
    def ue_trajs(self):
        """Gets the ue_trajs of this UeMobilityCollection.


        :return: The ue_trajs of this UeMobilityCollection.
        :rtype: List[UeTrajectoryCollection]
        """
        return self._ue_trajs

    @ue_trajs.setter
    def ue_trajs(self, ue_trajs):
        """Sets the ue_trajs of this UeMobilityCollection.


        :param ue_trajs: The ue_trajs of this UeMobilityCollection.
        :type ue_trajs: List[UeTrajectoryCollection]
        """
        if ue_trajs is None:
            raise ValueError("Invalid value for `ue_trajs`, must not be `None`")  # noqa: E501
        if ue_trajs is not None and len(ue_trajs) < 1:
            raise ValueError("Invalid value for `ue_trajs`, number of items must be greater than or equal to `1`")  # noqa: E501

        self._ue_trajs = ue_trajs
