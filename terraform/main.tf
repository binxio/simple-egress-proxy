//   Copyright 2021 binx.io B.V.
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//
resource "google_cloud_run_service" "push_subscription_proxy" {
  name     = "push-subscription-proxy"
  location = data.google_client_config.current.region

  template {
    spec {
      service_account_name = google_service_account.push_subscription_proxy.email
      containers {
        image = "gcr.io/binx-io-public/simple-egress-proxy:0.1.0"
        args  = ["--target-url", "https://httpbin.org/anything/event"]
      }
    }
  }
}

resource "google_service_account" "push_subscription_proxy" {
  account_id   = "push-subscription-proxy"
  display_name = "Pub/Sub push subscription proxy"
}

resource "google_cloud_run_service_iam_binding" "push_subscription_proxy_run_invokers" {
  location = google_cloud_run_service.push_subscription_proxy.location
  project  = google_cloud_run_service.push_subscription_proxy.project
  service  = google_cloud_run_service.push_subscription_proxy.name
  role     = "roles/run.invoker"
  members = [
    format("serviceAccount:%s", google_service_account.push_subscription_proxy.email)
  ]
  depends_on = [google_cloud_run_service.push_subscription_proxy]
}

resource "google_pubsub_topic" "notifications" {
  name = "notifications"
}

resource "google_pubsub_subscription" "proxied_push_subscription" {
  name  = "proxied-push-subscription"
  topic = google_pubsub_topic.notifications.name

  push_config {
    push_endpoint = google_cloud_run_service.push_subscription_proxy.status[0].url
    oidc_token {
      service_account_email = google_service_account.push_subscription_proxy.email
    }
  }
}
