// Emits the 7 "strong" hub mapping files for the next consolidation wave.
import fs from 'node:fs';
import path from 'node:path';
const ROOT = path.dirname(new URL(import.meta.url).pathname);
const SK = process.env.HOME + '/.claude/skills';

const families = {
  'chrome-extension': {
    family: 'chrome-extension',
    hubs: {
      'chrome-extension-expert': {
        keepExisting: false,
        title: 'Chrome Extension Development (MV3, APIs, packaging, security)',
        spokes: ['chrome-badge-metrics','chrome-dev','chrome-extension-packaging','chrome-extension-security-reviewer','chrome-identity-oauth','chrome-mv3-advanced','chrome-native-messaging','chrome-notifications-patterns','chrome-offscreen-documents','chrome-storage-patterns','chrome-tabs-management','extension-e2e-testing','extension-message-bridge','markdown-rendering-browser','mv3-service-worker-expert','shadow-dom-component-authoring','websocket-extension-patterns']
      }
    }
  },
  'ai-agent': {
    family: 'ai-agent',
    hubs: {
      'ai-agent-engineering': {
        keepExisting: false,
        title: 'AI & Agent Engineering (agents, MCP, LLMs, RAG, prompting)',
        spokes: ['a2a-interop','agent-council','agent-ecosystem','agent-harness-construction','agent-workflow-builder_ai_toolkit','ai-datastores','ai-languages','anthropic-sdk','autonomous-loops','continuous-learning-v2','llm-context-engineering','llm-integration-reviewer','llm-models','mcp-builder','mcp-servers','prompt-engineering','rag-architecture']
      }
    }
  },
  'software': {
    family: 'software',
    hubs: {
      'programming-languages': {
        keepExisting: false,
        title: 'Programming Languages (Python, Go, TypeScript, JS/Node, Kotlin/Compose)',
        spokes: ['python-patterns','go-patterns','javascript-nodejs','typescript-expert','typescript-advanced-types','compose-multiplatform-patterns','javascript-node-html-css-debugging-expert']
      },
      'software-engineering-patterns': {
        keepExisting: false,
        title: 'Software Engineering Patterns & Practices (architecture, APIs, debugging, reviews)',
        spokes: ['api-design-patterns','backend-patterns','microservices-patterns','express-patterns','coding-patterns','coding-standards','debugging','debugging-strategies','software-architect','code-reviewer','performance-profiling-expert','sse-streaming-patterns','job-scheduling-patterns','alarm-scheduler-patterns','auth-checker-patterns','diagnostic-registry-patterns','ops-registry-patterns','playbook-matcher-patterns','template-config-patterns','email-notification-patterns','web-auth-patterns','indexeddb-patterns','glean-llm-client-patterns','salesforce-scraping-patterns']
      }
    }
  },
  'integration': {
    family: 'integration',
    hubs: {
      'integration-clients': {
        keepExisting: false,
        title: 'Integration & API Clients (Jira, Monday, Slack, Salesforce, Glean, Aha)',
        spokes: ['aha-api','glean-dev','jira-developer-expert','jira-extension-client','monday-api','monday-dev','salesforce-developer-expert','slack-dev']
      }
    }
  },
  'devops': {
    family: 'devops',
    hubs: {
      'devops-infra': {
        keepExisting: false,
        title: 'DevOps, Infrastructure & Observability (Docker, K8s, CI/CD, Terraform, logging)',
        spokes: ['cicd-pipelines','docker-containers','kubernetes-networking','linux-sysadmin','nodejs-observability','pino-structured-logging','sentry-monitoring','shell-scripting','terraform-kafka-infra']
      }
    }
  },
  'frontend': {
    family: 'frontend',
    hubs: {
      'frontend-ui': {
        keepExisting: false,
        title: 'Frontend & UI/UX (design, HTML/CSS, web/mobile, accessibility reviews)',
        spokes: ['accessibility-ux-reviewer','frontend-design','frontend-design-ui-ux-expert','html-css','mobile-ios-design','ui-ux-pro-max','vanilla-js-ui-reviewer','web-design']
      }
    }
  }
};

let wrote = [];
for (const [key, m] of Object.entries(families)) {
  m.skillsRoot = SK;
  m.standalone = [];
  // validate spokes exist
  const missing = [];
  for (const def of Object.values(m.hubs)) for (const sp of def.spokes) if (!fs.existsSync(path.join(SK, sp, 'SKILL.md'))) missing.push(sp);
  const out = path.join(ROOT, `${key}-mapping.json`);
  fs.writeFileSync(out, JSON.stringify(m, null, 2));
  const total = Object.values(m.hubs).reduce((a, h) => a + h.spokes.length, 0);
  wrote.push(`${key}-mapping.json: ${Object.keys(m.hubs).length} hub(s), ${total} spokes${missing.length ? ' | MISSING: ' + missing.join(',') : ' ✓'}`);
}
console.log(wrote.join('\n'));
