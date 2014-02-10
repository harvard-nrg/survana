package dashboard

import (
        "neuroinformatics.harvard.edu/survana/auth"
       )

type Config struct {
    Authentication *auth.Config `json:"authentication,omitempty"`
    AllowRegistration bool      `json:"allow_registration,omitempty"`
}