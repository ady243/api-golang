import { App, TerraformStack } from 'cdktf';
import { AwsProvider } from '@cdktf/provider-aws/lib/provider';
import { SecurityGroup } from '@cdktf/provider-aws/lib/security-group';
import { Instance } from '@cdktf/provider-aws/lib/instance';
import { DbInstance } from '@cdktf/provider-aws/lib/db-instance';
import { DbSubnetGroup } from '@cdktf/provider-aws/lib/db-subnet-group';
import { ElasticacheCluster as CacheCluster } from '@cdktf/provider-aws/lib/elasticache-cluster';
import { DockerProvider } from '@cdktf/provider-docker/lib/provider';
import { Container as DockerContainer } from '@cdktf/provider-docker/lib/container';
import { Network } from '@cdktf/provider-docker/lib/network';

class MyStack extends TerraformStack {
  constructor(scope: App, id: string) {
    super(scope, id);

    // AWS Provider
    new AwsProvider(this, 'AWS', {
      region: 'us-east-1', 
    });

    // Docker Provider
    new DockerProvider(this, 'docker', {});

    // Security Group for EC2
    const securityGroup = new SecurityGroup(this, 'SecurityGroup', {
      name: 'my-unique-security-group-2',
      description: 'Allow inbound HTTP and SSH',
      ingress: [
        { fromPort: 22, toPort: 22, protocol: 'tcp', cidrBlocks: ['0.0.0.0/0'] },
        { fromPort: 80, toPort: 80, protocol: 'tcp', cidrBlocks: ['0.0.0.0/0'] },
        { fromPort: 3003, toPort: 3003, protocol: 'tcp', cidrBlocks: ['0.0.0.0/0'] },
      ],      
    });

    // EC2 Instance for hosting the Golang application
    new Instance(this, 'MyEC2Instance', {
      ami: 'ami-00c61e7eaac151f16',
      instanceType: 't2.micro',
      keyName: 'my-new-ec2-key',
      securityGroups: [securityGroup.name],
      tags: { Name: 'my-golang-api-instance' },
    });

    // Create DbSubnetGroup (for PostgreSQL)
    const dbSubnetGroup = new DbSubnetGroup(this, 'PostgresDbSubnetGroup', {
      subnetIds: [
        'subnet-0826c26bf860caa19',
        'subnet-09aec8d5d201e8ad5',
        'subnet-006a1b1bac9f65cf2',
      ],
      name: 'my-postgres-subnet-group-2',
      tags: { Name: 'PostgresSubnetGroup' },
    });

    // PostgreSQL via RDS
    new DbInstance(this, 'PostgresDB', {
      engine: 'postgres',
      instanceClass: 'db.t3.micro',
      allocatedStorage: 20,
      username: 'postgres',
      password: process.env.POSTGRES_PASSWORD || 'defaultPassword',
      publiclyAccessible: true,
      vpcSecurityGroupIds: [securityGroup.id],
      dbSubnetGroupName: dbSubnetGroup.id,
    });

    // Redis via ElastiCache
    new CacheCluster(this, 'RedisCluster', {
      clusterId: 'my-unique-redis-cluster-2',
      engine: 'redis',
      nodeType: 'cache.t3.micro',
      numCacheNodes: 1,
      securityGroupIds: [securityGroup.id],
    });

    // Create Docker Network
    const network = new Network(this, 'network', {
      name: 'my-network',
      driver: 'bridge',
    });

    // Docker Container for PostgreSQL
    new DockerContainer(this, 'postgres', {
      image: 'postgres:latest',
      name: 'postgres',
      ports: [
        { internal: 5432, external: 5433 }, 
      ],
      env: [
        `POSTGRES_USER=${process.env.POSTGRES_USER || 'postgres'}`,
        `POSTGRES_PASSWORD=${process.env.POSTGRES_PASSWORD || 'defaultPassword'}`,
        `POSTGRES_DB=${process.env.POSTGRES_DATABASE || 'postgres'}`,
      ],
      networksAdvanced: [
        {
          name: network.name,
          aliases: ['postgres'],
        },
      ],
    });

    // Docker Container for Redis
    new DockerContainer(this, 'redis', {
      image: 'redis:latest',
      name: 'redis',
      ports: [
        { internal: 6379, external: 6381 }, 
      ],
      networksAdvanced: [
        {
          name: network.name,
          aliases: ['redis'],
        },
      ],
    });

    // Docker Container for Dragonfly
    new DockerContainer(this, 'dragonfly', {
      image: 'docker.dragonflydb.io/dragonflydb/dragonfly',
      name: 'dragonfly',
      ports: [
        { internal: 6379, external: 6382 }, // Utiliser un port externe différent
      ],
      networksAdvanced: [
        {
          name: network.name,
          aliases: ['dragonfly'],
        },
      ],
    });

    // Docker Container for the Golang application
    new DockerContainer(this, 'api', {
      image: 'my-golang-app',
      name: 'api-unique',
      ports: [
        { internal: 3003, external: 3004 },
      ],
      restart: 'on-failure',
      env: [
        `POSTGRES_HOST=postgres`,
        `POSTGRES_USER=postgres`,
        `POSTGRES_PASSWORD=password`,
        `POSTGRES_DB=postgres`,
        `REDIS_ADDR=redis:6379`,
        `GOOGLE_CLIENT_ID=${process.env.GOOGLE_CLIENT_ID || '352689561996-btrtmukipkgudca8dr10jb7m7j73knbi.apps.googleusercontent.com'}`,
        `GOOGLE_CLIENT_SECRET=${process.env.GOOGLE_CLIENT_SECRET || 'Q6J9J6Q6J9J6Q6J9J6Q6J6Q6'}`,
        `GOOGLE_REDIRECT_URI=${process.env.GOOGLE_REDIRECT_URI || 'com.googleusercontent.apps.352689561996-btrtmukipkgudca8dr10jb7m7j73knbi:/callback'}`,
        `OPENAI_API_KEY=${process.env.OPENAI_API_KEY || 'sk-V9AhdR3zhHX2GW2X9s8lJmIpH2cV6hPRPULaoRPTfTT3BlbkFJqAJiCDlothHFm25lad6BDdafGb7y_wJcZT-KaiAKEA'}`,
        `DRAGONFLY_HOST=dragonfly`,
        `DRAGONFLY_PORT=6379`,
      ],
      networksAdvanced: [
        {
          name: network.name,
          aliases: ['api'],
        },
      ],
    });
  }
}

const app = new App();
new MyStack(app, 'my-stack');
app.synth();