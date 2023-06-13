package webhook

import (
	"encoding/json"
	"fmt"
	"kubearchvalidator/pkg/registry"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AdmissionController struct {
	Server *http.Server
}

func NewAdmissionController(addr string) *AdmissionController {
	mux := http.NewServeMux()
	mux.HandleFunc("/admit", func(w http.ResponseWriter, r *http.Request) {
		var review admissionv1.AdmissionReview
		if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
			http.Error(w, fmt.Sprintf("could not decode body: %v", err), http.StatusBadRequest)
			return
		}

		response := Admit(review)
		response.UID = review.Request.UID

		res, err := json.Marshal(admissionv1.AdmissionReview{
			Response: response,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(res); err != nil {
			http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
			return
		}
	})

	return &AdmissionController{
		Server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

func Admit(review admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	var pod corev1.Pod
	if err := json.Unmarshal(review.Request.Object.Raw, &pod); err != nil {
		return &admissionv1.AdmissionResponse{
			Allowed: false,
			Result: &metav1.Status{
				Message: fmt.Sprintf("could not unmarshal pod: %v", err),
			},
		}
	}

	for _, container := range pod.Spec.Containers {
		image := container.Image
		if !registry.CheckImageArchitecture(image) {
			return &admissionv1.AdmissionResponse{
				Allowed: false,
				Result: &metav1.Status{
					Message: fmt.Sprintf("image %s does not match node architecture", image),
				},
			}
		}
	}

	return &admissionv1.AdmissionResponse{
		Allowed: true,
	}
}
