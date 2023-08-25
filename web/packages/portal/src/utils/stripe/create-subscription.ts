import { getAvailableProducts } from "@/utils/stripe/available-products";
import { getStripeClient } from "@/utils/stripe/stripe-client";

export const createSubscription = async (
  customerId: string,
  workspaceId: string,
  chargeable: boolean
) => {
  const stripeClient = getStripeClient(chargeable);
  return await stripeClient.subscriptions.create({
    customer: customerId,
    metadata: {
      workspace_id: workspaceId,
    },
    billing_thresholds: {
      amount_gte: 50,
      reset_billing_cycle_anchor: false,
    },
    items: getAvailableProducts(chargeable),
  });
};
