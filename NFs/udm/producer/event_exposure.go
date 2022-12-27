package producer

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/free5gc/http_wrapper"
	"github.com/free5gc/openapi/models"
	udm_context "github.com/free5gc/udm/context"
	"github.com/free5gc/udm/logger"
)

func HandleCreateEeSubscription(request *http_wrapper.Request) *http_wrapper.Response {
	logger.EeLog.Infoln("Handle Create EE Subscription")

	eesubscription := request.Body.(models.EeSubscription)
	ueIdentity := request.Params["ueIdentity"]

	createdEESubscription, problemDetails := CreateEeSubscriptionProcedure(ueIdentity, eesubscription)
	if createdEESubscription != nil {
		return http_wrapper.NewResponse(http.StatusCreated, nil, createdEESubscription)
	} else if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	} else {
		problemDetails = &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "UNSPECIFIED_NF_FAILURE",
		}
		return http_wrapper.NewResponse(http.StatusInternalServerError, nil, problemDetails)
	}
}

// TODO: complete this procedure based on TS 29503 5.5
func CreateEeSubscriptionProcedure(ueIdentity string,
	eesubscription models.EeSubscription) (*models.CreatedEeSubscription, *models.ProblemDetails) {
	udmSelf := udm_context.UDM_Self()

	logger.EeLog.Debugf("udIdentity: %s", ueIdentity)
	switch {
	// GPSI (MSISDN identifier) represents a single UE
	case strings.HasPrefix(ueIdentity, "msisdn-"):
		fallthrough
	// GPSI (External identifier) represents a single UE
	case strings.HasPrefix(ueIdentity, "extid-"):
		if ue, ok := udmSelf.UdmUeFindByGpsi(ueIdentity); ok {
			id, err := udmSelf.EeSubscriptionIDGenerator.Allocate()
			if err != nil {
				problemDetails := &models.ProblemDetails{
					Status: http.StatusInternalServerError,
					Cause:  "UNSPECIFIED_NF_FAILURE",
				}
				return nil, problemDetails
			}

			subscriptionID := strconv.Itoa(int(id))
			ue.EeSubscriptions[subscriptionID] = &eesubscription
			createdEeSubscription := &models.CreatedEeSubscription{
				EeSubscription: &eesubscription,
			}
			return createdEeSubscription, nil
		} else {
			problemDetails := &models.ProblemDetails{
				Status: http.StatusNotFound,
				Cause:  "USER_NOT_FOUND",
			}
			return nil, problemDetails
		}
	// external groupID represents a group of UEs
	case strings.HasPrefix(ueIdentity, "extgroupid-"):
		id, err := udmSelf.EeSubscriptionIDGenerator.Allocate()
		if err != nil {
			problemDetails := &models.ProblemDetails{
				Status: http.StatusInternalServerError,
				Cause:  "UNSPECIFIED_NF_FAILURE",
			}
			return nil, problemDetails
		}
		subscriptionID := strconv.Itoa(int(id))
		createdEeSubscription := &models.CreatedEeSubscription{
			EeSubscription: &eesubscription,
		}

		udmSelf.UdmUePool.Range(func(key, value interface{}) bool {
			ue := value.(*udm_context.UdmUeContext)
			if ue.ExternalGroupID == ueIdentity {
				ue.EeSubscriptions[subscriptionID] = &eesubscription
			}
			return true
		})
		return createdEeSubscription, nil
	// represents any UEs
	case ueIdentity == "anyUE":
		id, err := udmSelf.EeSubscriptionIDGenerator.Allocate()
		if err != nil {
			problemDetails := &models.ProblemDetails{
				Status: http.StatusInternalServerError,
				Cause:  "UNSPECIFIED_NF_FAILURE",
			}
			return nil, problemDetails
		}
		subscriptionID := strconv.Itoa(int(id))
		createdEeSubscription := &models.CreatedEeSubscription{
			EeSubscription: &eesubscription,
		}
		udmSelf.UdmUePool.Range(func(key, value interface{}) bool {
			ue := value.(*udm_context.UdmUeContext)
			ue.EeSubscriptions[subscriptionID] = &eesubscription
			return true
		})
		return createdEeSubscription, nil
	default:
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "MANDATORY_IE_INCORRECT",
			InvalidParams: []models.InvalidParam{
				{
					Param:  "ueIdentity",
					Reason: "incorrect format",
				},
			},
		}
		return nil, problemDetails
	}
}

func HandleDeleteEeSubscription(request *http_wrapper.Request) *http_wrapper.Response {
	ueIdentity := request.Params["ueIdentity"]
	subscriptionID := request.Params["subscriptionID"]

	DeleteEeSubscriptionProcedure(ueIdentity, subscriptionID)
	return http_wrapper.NewResponse(http.StatusNoContent, nil, nil)
}

// TODO: complete this procedure based on TS 29503 5.5
func DeleteEeSubscriptionProcedure(ueIdentity string, subscriptionID string) {
	udmSelf := udm_context.UDM_Self()

	switch {
	case strings.HasPrefix(ueIdentity, "msisdn-"):
		fallthrough
	case strings.HasPrefix(ueIdentity, "extid-"):
		if ue, ok := udmSelf.UdmUeFindByGpsi(ueIdentity); ok {
			delete(ue.EeSubscriptions, subscriptionID)
		}
	case strings.HasPrefix(ueIdentity, "extgroupid-"):
		udmSelf.UdmUePool.Range(func(key, value interface{}) bool {
			ue := value.(*udm_context.UdmUeContext)
			if ue.ExternalGroupID == ueIdentity {
				delete(ue.EeSubscriptions, subscriptionID)
			}
			return true
		})
	case ueIdentity == "anyUE":
		udmSelf.UdmUePool.Range(func(key, value interface{}) bool {
			ue := value.(*udm_context.UdmUeContext)
			delete(ue.EeSubscriptions, subscriptionID)
			return true
		})
	}
	if id, err := strconv.ParseInt(subscriptionID, 10, 64); err != nil {
		logger.EeLog.Warnf("subscriptionID covert type error: %+v", err)
	} else {
		udmSelf.EeSubscriptionIDGenerator.FreeID(id)
	}
}

func HandleUpdateEeSubscription(request *http_wrapper.Request) *http_wrapper.Response {
	logger.EeLog.Infoln("Handle Update EE subscription")
	logger.EeLog.Warnln("Update EE Subscription is not implemented")

	patchList := request.Body.([]models.PatchItem)
	ueIdentity := request.Params["ueIdentity"]
	subscriptionID := request.Params["subscriptionID"]

	problemDetails := UpdateEeSubscriptionProcedure(ueIdentity, subscriptionID, patchList)
	if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	} else {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, nil)
	}
}

// TODO: complete this procedure based on TS 29503 5.5
func UpdateEeSubscriptionProcedure(ueIdentity string, subscriptionID string,
	patchList []models.PatchItem) *models.ProblemDetails {
	udmSelf := udm_context.UDM_Self()

	switch {
	case strings.HasPrefix(ueIdentity, "msisdn-"):
		fallthrough
	case strings.HasPrefix(ueIdentity, "extid-"):
		if ue, ok := udmSelf.UdmUeFindByGpsi(ueIdentity); ok {
			if _, ok := ue.EeSubscriptions[subscriptionID]; ok {
				for _, patchItem := range patchList {
					logger.EeLog.Debugf("patch item: %+v", patchItem)
					// TODO: patch the Eesubscription
				}
				return nil
			} else {
				problemDetails := &models.ProblemDetails{
					Status: http.StatusNotFound,
					Cause:  "SUBSCRIPTION_NOT_FOUND",
				}
				return problemDetails
			}
		} else {
			problemDetails := &models.ProblemDetails{
				Status: http.StatusNotFound,
				Cause:  "SUBSCRIPTION_NOT_FOUND",
			}
			return problemDetails
		}
	case strings.HasPrefix(ueIdentity, "extgroupid-"):
		udmSelf.UdmUePool.Range(func(key, value interface{}) bool {
			ue := value.(*udm_context.UdmUeContext)
			if ue.ExternalGroupID == ueIdentity {
				if _, ok := ue.EeSubscriptions[subscriptionID]; ok {
					for _, patchItem := range patchList {
						logger.EeLog.Debugf("patch item: %+v", patchItem)
						// TODO: patch the Eesubscription
					}
				}
			}
			return true
		})
		return nil
	case ueIdentity == "anyUE":
		udmSelf.UdmUePool.Range(func(key, value interface{}) bool {
			ue := value.(*udm_context.UdmUeContext)
			if _, ok := ue.EeSubscriptions[subscriptionID]; ok {
				for _, patchItem := range patchList {
					logger.EeLog.Debugf("patch item: %+v", patchItem)
					// TODO: patch the Eesubscription
				}
			}
			return true
		})
		return nil
	default:
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "MANDATORY_IE_INCORRECT",
			InvalidParams: []models.InvalidParam{
				{
					Param:  "ueIdentity",
					Reason: "incorrect format",
				},
			},
		}
		return problemDetails
	}
}
