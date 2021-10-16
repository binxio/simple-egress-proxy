simple egress proxy
==================
In projects protected by a VPC service control perimeter, new push subscriptions cannot be created unless the push 
endpoints are set to Cloud Run services with default run.app URLs (custom domains won't work). If you try
to create a push subscription, the client get the error: `Request is prohibited by organization's policy.`

This simple proxy can be used to work around that problem. You deploy this image as Cloud Run service and subscribe the
cloud run service to the topic. This egress proxy will forward the event to the real destination.

```hcl
resource "google_cloud_run_service" "push_subscription_proxy" {
  name     = "push-subscription-proxy"
  location = "europe-west4"

  template {
    spec {
      service_account_name = google_service_account.push_subscription_proxy.email
      containers {
        image = "gcr.io/binx-io-public/simple-egress-proxy:0.1.0"
        args = [ "--target-url", "https://httpbin.org/anything/event"]
      }
    }
  }
}
```

Now you can create a Google Pub/Sub subscription:

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

See also [How to configure a Google Pub/Sub push subscription within a VPC Service control perimeter](https://binx.io/blog/2021/10/16/how-to-configure…ontrol-perimeter)