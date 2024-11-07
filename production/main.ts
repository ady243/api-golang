import { App, TerraformStack } from "cdktf";
import { Construct } from "constructs";
import { GoogleProvider } from "@cdktf/provider-google/lib/provider";
import { CloudRunService } from "@cdktf/provider-google/lib/cloud-run-service";
import { CloudRunServiceIamMember } from "@cdktf/provider-google/lib/cloud-run-service-iam-member";
import * as dotenv from 'dotenv';

// Charger les variables d'environnement depuis le fichier .env
dotenv.config();

class MyStack extends TerraformStack {
  constructor(scope: Construct, id: string) {
    super(scope, id);

    // Configure le provider Google
    new GoogleProvider(this, "Google", {
      project: process.env.GCP_PROJECT_ID,
      region: "us-central1",
    });

    // Définir le service Cloud Run
    const cloudRunService = new CloudRunService(this, "MyCloudRunService", {
      name: "my-cloud-run-service",
      location: "us-central1",
      template: {
        spec: {
          containers: [
            {
              image: `gcr.io/${process.env.GCP_PROJECT_ID}/your-docker-image:latest`,
              ports: [
                {
                  containerPort: 3003,
                },
              ],
              env: [
                { name: "POSTGRES_HOST", value: process.env.POSTGRES_HOST },
                { name: "POSTGRES_USER", value: process.env.POSTGRES_USER },
                { name: "POSTGRES_PASSWORD", value: process.env.POSTGRES_PASSWORD },
                { name: "POSTGRES_DB", value: process.env.POSTGRES_DB },
                { name: "DRAGONFLY_HOST", value: process.env.DRAGONFLY_HOST },
                { name: "GOOGLE_CLIENT_SECRET", value: process.env.GOOGLE_CLIENT_SECRET },
                { name: "GOOGLE_REDIRECT_URI", value: process.env.GOOGLE_REDIRECT_URI },
                { name: "OPENAI_API_KEY", value: process.env.OPENAI_API_KEY },
              ],
            },
          ],
        },
      },
    });

    // Configurer les permissions IAM pour permettre l'accès public
    new CloudRunServiceIamMember(this, "MyCloudRunServiceIamMember", {
      service: cloudRunService.name,
      location: cloudRunService.location,
      role: "roles/run.invoker",
      member: "allUsers",
    });
  }
}

const app = new App();
new MyStack(app, "api-golang");
app.synth();