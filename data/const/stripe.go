package data

const (
	StripeEventCheckoutSessionCompleted             = "checkout.session.completed"
	StripeEventCheckoutSessionAsyncPaymentSucceeded = "checkout.session.async_payment_succeeded"
	StripeEventCheckoutSessionExpired               = "checkout.session.expired"
	StripeEventCustomerSourceExpiring               = "customer.source.expiring"
	StripeEventCustomerSubscriptionDeleted          = "customer.subscription.deleted"
	StripeEventCustomerSubscriptionUpdated          = "customer.subscription.updated"
	StripeEventRadarEarlyFraudWarningCreated        = "radar.early_fraud_warning.created"
	StripeEventInvoicePaymentActionRequired         = "invoice.payment_action_required"
	StripeEventCustomerSubscriptionTrialWillEnd     = "customer.subscription.trial_will_end"
	StripeEventInvoicePaymentFailed                 = "invoice.payment_failed"
)
