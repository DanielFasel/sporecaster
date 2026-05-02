async function render() {
  const spore = await fetch('/api/spore').then(r => r.json());

  document.title = spore.app;
  document.getElementById('app-name').textContent = spore.app;

  const langBadge = document.getElementById('app-lang');
  langBadge.textContent = spore.language;
  langBadge.classList.add('badge', 'lang-badge');

  document.getElementById('app-desc').textContent = spore.description || '';

  // build parent → children map
  const byParent = {};
  for (const pkg of spore.packages || []) {
    const key = pkg.parent || '';
    if (!byParent[key]) byParent[key] = [];
    byParent[key].push(pkg);
  }

  const tree = document.getElementById('tree');

  if (spore.core) {
    tree.appendChild(makeCard({ ...spore.core, isCore: true }, byParent));
  }

  for (const pkg of byParent[''] || []) {
    tree.appendChild(makeCard(pkg, byParent));
  }

  // channels
  const channelsEl = document.getElementById('channels');
  for (const ch of spore.channels || []) {
    channelsEl.appendChild(makeChannelCard(ch));
  }
}

function makeCard(pkg, byParent) {
  const card = document.createElement('div');
  card.className = 'card' + (pkg.isCore ? ' core' : '');

  // ── header ──────────────────────────────
  const header = document.createElement('div');
  header.className = 'card-header';

  const name = document.createElement('span');
  name.className = 'pkg-name';
  name.textContent = pkg.isCore ? pkg.name + ' (core)' : pkg.name;
  header.appendChild(name);

  if (pkg.kind) {
    const kind = document.createElement('span');
    kind.className = 'badge kind-' + pkg.kind;
    kind.textContent = pkg.kind;
    header.appendChild(kind);
  }

  const hasExports = pkg.exports && pkg.exports.length > 0;
  if (hasExports) {
    const toggle = document.createElement('span');
    toggle.className = 'exports-toggle';
    toggle.textContent = '▸';
    header.appendChild(toggle);
  }

  card.appendChild(header);

  // ── description ─────────────────────────
  if (pkg.description) {
    const desc = document.createElement('p');
    desc.className = 'pkg-desc';
    desc.textContent = pkg.description.trim();
    card.appendChild(desc);
  }

  // ── imports ──────────────────────────────
  if (pkg.imports && pkg.imports.length > 0) {
    const row = document.createElement('div');
    row.className = 'imports';
    const label = document.createElement('span');
    label.className = 'imports-label';
    label.textContent = 'imports:';
    row.appendChild(label);
    for (const imp of pkg.imports) {
      const badge = document.createElement('span');
      badge.className = 'import-badge';
      badge.textContent = imp;
      row.appendChild(badge);
    }
    card.appendChild(row);
  }

  // ── exports (collapsed by default) ───────
  if (hasExports) {
    const section = document.createElement('div');
    section.className = 'exports-section';

    const expLabel = document.createElement('div');
    expLabel.className = 'exports-label';
    expLabel.textContent = 'exports';
    section.appendChild(expLabel);

    for (const exp of pkg.exports) {
      const row = document.createElement('div');
      row.className = 'export-row';
      const expName = document.createElement('span');
      expName.className = 'export-name';
      expName.textContent = exp.name;
      const expKind = document.createElement('span');
      expKind.className = 'badge export-kind-' + exp.kind;
      expKind.textContent = exp.kind;
      row.appendChild(expName);
      row.appendChild(expKind);
      section.appendChild(row);
    }

    card.appendChild(section);

    header.style.cursor = 'pointer';
    header.addEventListener('click', () => {
      const open = section.classList.toggle('open');
      header.querySelector('.exports-toggle').textContent = open ? '▾' : '▸';
    });
  }

  // ── sub-packages ─────────────────────────
  const children = byParent[pkg.name];
  if (children && children.length > 0) {
    const group = document.createElement('div');
    group.className = 'children';
    for (const child of children) {
      group.appendChild(makeCard(child, byParent));
    }
    card.appendChild(group);
  }

  return card;
}

function makeChannelCard(ch) {
  const card = document.createElement('div');
  card.className = 'channel-card';

  const header = document.createElement('div');
  header.className = 'channel-header';

  const name = document.createElement('span');
  name.className = 'channel-name';
  name.textContent = ch.name;
  header.appendChild(name);

  const type = document.createElement('span');
  type.className = 'badge type-badge';
  type.textContent = ch.type;
  header.appendChild(type);

  card.appendChild(header);

  if (ch.commands && ch.commands.length > 0) {
    const list = document.createElement('div');
    list.className = 'commands';
    for (const cmd of ch.commands) {
      const row = document.createElement('div');
      row.className = 'command';
      const usage = document.createElement('div');
      usage.className = 'command-usage';
      usage.textContent = cmd.usage || cmd.name;
      row.appendChild(usage);
      if (cmd.description) {
        const desc = document.createElement('div');
        desc.className = 'command-desc';
        desc.textContent = cmd.description.trim();
        row.appendChild(desc);
      }
      list.appendChild(row);
    }
    card.appendChild(list);
  }

  return card;
}

render().catch(err => {
  document.getElementById('tree').textContent = 'Error loading spore: ' + err.message;
});
