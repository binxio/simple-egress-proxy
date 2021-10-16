simple egress proxy
==================
In projects protected by a VPC service control perimeter, new push subscriptions cannot be created unless the push 
endpoints are set to Cloud Run services with default run.app URLs (custom domains won't work). 

This simple proxy can be used to work around that problem. You deploy this image as Cloud Run service and subscribe the
cloud run service to the topic. This egress will proxy to the real target.

```hcl
resource "google_cloud_run_service" "push_subscription_proxy" {
  name     = "push-subscription-proxy"
  location = "europe-west4"

  template {
    spec {
      service_account_name = google_service_account.push_subscription_proxy.email
      containers {
        image = "gcr.io/binx-io-public/simple-egress-proxy:latest"
        args = [ "--target-url", "https://httpbin.org/anything"]
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
}

```

```hcl
resource "google_pubsub_topic" "notifications" {
  name    = "notifications"
}

resource "google_pubsub_subscription" "proxied_push_subscription" {
  name    = "proxied-push-subscription"
  topic   = google_pubsub_topic.notifications.name

  push_config {
    push_endpoint = google_cloud_run_service.push_subscription_proxy.status[0].url
    oidc_token {
      service_account_email = google_service_account.push_subscription_proxy.email
    }
  }
}
```
