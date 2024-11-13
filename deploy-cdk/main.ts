import { App, TerraformStack } from 'cdktf';
import { AwsProvider } from '@cdktf/provider-aws/lib/provider';
import { SecurityGroup } from '@cdktf/provider-aws/lib/security-group';
import { Instance } from '@cdktf/provider-aws/lib/instance';
import { DbInstance } from '@cdktf/provider-aws/lib/db-instance';
import { ElasticacheCluster as CacheCluster } from '@cdktf/provider-aws/lib/elasticache-cluster';
import { DockerProvider } from '@cdktf/provider-docker/lib/provider';
import { Container as DockerContainer } from '@cdktf/provider-docker/lib/container';

class MyStack extends TerraformStack {
  constructor(scope: App, id: string) {
    super(scope, id);

    // AWS Provider
    new AwsProvider(this, 'AWS', {
      region: 'us-west-2',
    });

    // Docker Provider
    new DockerProvider(this, 'docker', {});

    // Security Group for EC2
    const securityGroup = new SecurityGroup(this, 'SecurityGroup', {
      name: 'my-security-group',
      description: 'Allow inbound HTTP and SSH',
      ingress: [
        { fromPort: 22, toPort: 22, protocol: 'tcp', cidrBlocks: ['0.0.0.0/0'] },
        { fromPort: 80, toPort: 80, protocol: 'tcp', cidrBlocks: ['0.0.0.0/0'] },
        { fromPort: 3003, toPort: 3003, protocol: 'tcp', cidrBlocks: ['0.0.0.0/0'] },
      ],
    });

    // EC2 Instance for hosting the Golang application
    new Instance(this, 'MyEC2Instance', {
      ami: 'ami-0984f4b9e98be44bf', 
      instanceType: 't2.micro',
      keyName: 'your-key-name',
      securityGroups: [securityGroup.name],
      tags: { Name: 'my-golang-api-instance' },
    });

    // PostgreSQL via RDS
    new DbInstance(this, 'PostgresDB', {
      engine: 'postgres',
      instanceClass: 'db.t3.micro',
      allocatedStorage: 20,
      username: 'postgres',
      password: 'postgrespassword',
      publiclyAccessible: true,
      vpcSecurityGroupIds: [securityGroup.id],
    });

    // Redis via ElastiCache
    new CacheCluster(this, 'RedisCluster', {
      clusterId: 'my-redis-cluster',
      engine: 'redis',
      nodeType: 'cache.t3.micro',
      numCacheNodes: 1,
      securityGroupIds: [securityGroup.id],
    });

    // Docker Container for the Golang application
    new DockerContainer(this, 'api', {
      image: 'my-golang-app',
      name: 'api',
      ports: [
        { internal: 3003, external: 3003 },
      ],
      env: [
        'POSTGRES_HOST=postgres',
        'POSTGRES_USER=postgres',
        'POSTGRES_PASSWORD=postgrespassword',
        'POSTGRES_DB=postgres',
        'REDIS_ADDR=redis:6379',
        'GOOGLE_CLIENT_ID=352689561996-btrtmukipkgudca8dr10jb7m7j73knbi.apps.googleusercontent.com',
        'GOOGLE_CLIENT_SECRET=Q6J9J6Q6J9J6Q6J9J6Q6J6Q6',
        'GOOGLE_REDIRECT_URI=com.googleusercontent.apps.352689561996-btrtmukipkgudca8dr10jb7m7j73knbi:/callback',
        'OPENAI_API_KEY=sk-V9AhdR3zhHX2GW2X9s8lJmIpH2cV6hPRPULaoRPTfTT3BlbkFJqAJiCDlothHFm25lad6BDdafGb7y_wJcZT-KaiAKEA',
        'DRAGONFLY_HOST=dragonfly',
        'DRAGONFLY_PORT=6379',
      ],
    });
  }
}

const app = new App();
new MyStack(app, 'my-stack');
app.synth();