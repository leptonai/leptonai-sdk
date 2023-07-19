import { cors } from "@/utils/cors";
import { stripeClient } from "@/utils/stripe";
import { getWorkspaceById } from "@/utils/workspace";
import { NextApiRequest, NextApiResponse } from "next";

async function handler(req: NextApiRequest, res: NextApiResponse) {
  try {
    const workspaceId = JSON.parse(req.body).workspace_id;
    const workspace = await getWorkspaceById(workspaceId, req.cookies);
    if (!workspace) {
      return res.status(401).send("You are not authorized to call this API");
    } else {
      const consumerId = workspace.consumer_id;
      const session = await stripeClient.billingPortal.sessions.create({
        customer: consumerId,
        return_url: "https://dashboard.lepton.ai/settings/usage",
      });
      res.status(200).json({ url: session.url });
    }
  } catch (err) {
    const errorMessage =
      err instanceof Error ? err.message : "Internal server error";
    res.status(500).send(`Error: ${errorMessage}`);
  }
}

export default cors(handler);
