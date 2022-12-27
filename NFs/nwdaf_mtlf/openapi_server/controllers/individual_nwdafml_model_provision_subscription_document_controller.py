import connexion
import six
from typing import Dict
from typing import Tuple
from typing import Union
from connexion.operations.openapi import OpenAPIOperation
from openapi_server.models.ml_model_addr import MLModelAddr
from openapi_server.models.nwdaf_ml_model_prov_subsc import NwdafMLModelProvSubsc  # noqa: E501
from openapi_server.models.problem_details import ProblemDetails  # noqa: E501
from openapi_server.models.redirect_response import RedirectResponse  # noqa: E501
from openapi_server import util

from openapi_server.models.nwdaf_ml_model_prov_notif import NwdafMLModelProvNotif
from openapi_server.models.ml_event_notif import MLEventNotif

global service_req

def delete_nwdafml_model_provision_subcription(subscription_id):  # noqa: E501
    """Delete an existing Individual NWDAF ML Model Provision Subscription.

     # noqa: E501

    :param subscription_id: String identifying a subscription to the Nnwdaf_MLModelProvision Service.
    :type subscription_id: str

    :rtype: Union[None, Tuple[None, int], Tuple[None, int, Dict[str, str]]
    """
    return 'do some magic!'


def update_nwdafml_model_provision_subcription(subscription_id=None, nwdaf_ml_model_prov_subsc=None):  # noqa: E501
    """update an existing Individual NWDAF ML Model Provision Subscription

     # noqa: E501

    :param subscription_id: String identifying a subscription to the Nnwdaf_MLModelProvision Service.
    :type subscription_id: str
    :param nwdaf_ml_model_prov_subsc: 
    :type nwdaf_ml_model_prov_subsc: dict | bytes

    :rtype: Union[NwdafMLModelProvSubsc, Tuple[NwdafMLModelProvSubsc, int], Tuple[NwdafMLModelProvSubsc, int, Dict[str, str]]
    """
    global service_req
    MLmodelRequest = connexion.operations.openapi.tmpbody
    ml_event = MLmodelRequest["mLEventSubscs"][0]["mLEvent"]

    if ml_event == "ML_model_provision_request_delay":
        print("ML model provision requested")
        print("Provisioning ML model...")
        return MLEventNotif(m_l_file_addr= MLModelAddr()).to_dict()

    if ml_event =="ML_model_provision_request_drop":
        print("ML model provision requested")
        print("Provisioning ML model...")
        return MLEventNotif(m_l_file_addr= MLModelAddr()).to_dict()

    if ml_event == "EXAMPLE":
        return NwdafMLModelProvSubsc(m_l_event_notifs=MLEventNotif().to_dict()).to_dict()
