import connexion
import six
from typing import Dict
from typing import Tuple
from typing import Union
import datetime

from openapi_server.models.analytics_data import AnalyticsData  # noqa: E501
from openapi_server.models.event_id_any_of import EventIdAnyOf  # noqa: E501
from openapi_server.models.event_id import EventId 
from openapi_server.models.event_filter import EventFilter  # noqa: E501
from openapi_server.models.event_reporting_requirement import EventReportingRequirement
from openapi_server.models.network_area_info import NetworkAreaInfo  # noqa: E501
from openapi_server.models.problem_details import ProblemDetails  # noqa: E501
from openapi_server.models.problem_details_analytics_info_request import ProblemDetailsAnalyticsInfoRequest  # noqa: E501
from openapi_server.models.target_ue_information import TargetUeInformation  # noqa: E501
from openapi_server import util
from openapi_server.models.redundant_transmission_exp_info import RedundantTransmissionExpInfo
from openapi_server.models.redundant_transmission_exp_per_ts import RedundantTransmissionExpPerTS
from openapi_server.models.redundant_transmission_exp_req import RedundantTransmissionExpReq
from openapi_server.models.redundant_transmission_exp_predict import RedTransExpPredict

service = 0

def get_nwdaf_analytics(event_id, ana_req=None, event_filter=None, supported_features=None, tgt_ue=None):  # noqa: E501

    service = connexion.operations.openapi.tmpqry
    event_id = service["event-id"]

    if event_id == "EXAMPLE":
        return AnalyticsData().to_dict()

    if event_id == "REDUNDANT_TRANSMISSION":

        return AnalyticsData(red_trans_infos=[RedundantTransmissionExpInfo(dnn="mnc", spatial_valid_con = NetworkAreaInfo(g_ran_node_ids="UERANSIM"), red_trans_exps=RedundantTransmissionExpPerTS(ts_start=datetime.datetime.now(), ts_duration=40, red_trans_exp= RedTransExpPredict(), ue_ratio=1,confidence=None))])


    else:
        return "Not Supported Service"