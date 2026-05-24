(() => {
  const urlParams = new URLSearchParams(window.location.search);
  const forceMock = urlParams.get('mock') === '1' || window.PVC_EXPORTER_AGENT_MOCK === true;
  const disableMock = urlParams.get('mock') === '0';
  const enableAutoMock = urlParams.get('mockAuto') === '1' || window.PVC_EXPORTER_AGENT_AUTO_MOCK === true;
  const allowAutoMock = !disableMock && enableAutoMock;
  let mockEnabled = forceMock;

  const mockFS = {
    '': [
      { name: 'reports', isDir: true, modTime: '2026-05-20T08:15:00Z' },
      { name: 'notes.txt', isDir: false, size: 1536, modTime: '2026-05-20T07:58:00Z' },
      { name: 'README.md', isDir: false, size: 4096, modTime: '2026-05-20T07:42:00Z' },
    ],
    reports: [
      { name: 'daily.csv', isDir: false, size: 81234, modTime: '2026-05-20T06:30:00Z' },
      { name: 'weekly.csv', isDir: false, size: 214567, modTime: '2026-05-19T21:12:00Z' },
    ],
  };

  function mockResponse(payload, status = 200) {
    return Promise.resolve({
      ok: status >= 200 && status < 300,
      status,
      json: () => Promise.resolve(payload),
    });
  }

  function normalizePath(path) {
    return (path || '').replace(/^\/+|\/+$/g, '');
  }

  function mockList(path) {
    const p = normalizePath(path);
    if (!(p in mockFS)) return mockResponse({ entries: [] }, 200);
    return mockResponse({ entries: mockFS[p] }, 200);
  }

  function mockConfig() {
    return mockResponse({
      readonly: true,
      forceRW: false,
      pvcWatch: true,
      cluster: 'mock-cluster',
      namespace: 'demo',
      pvc: 'mock-pvc',
      pod: 'mock-pod',
    });
  }

  function mockSpace() {
    return mockResponse({ used: 3.2 * 1024 * 1024 * 1024, total: 10 * 1024 * 1024 * 1024 });
  }

  function mockNotSupported() {
    return Promise.reject(new Error('Action is disabled in mock mode'));
  }

  const API = {
    list: (path) => mockEnabled ? mockList(path) : fetch(`/api/files?path=${encodeURIComponent(path)}`),
    download: (path) => {
      if (mockEnabled) {
        window.alert(`Mock mode: download not available for ${path}`);
        return;
      }
      window.location = `/api/download?path=${encodeURIComponent(path)}`;
    },
    upload: (dir, file, onProgress) => {
      if (mockEnabled) return mockNotSupported();
      return new Promise((resolve, reject) => {
        const f = new FormData(); f.append('file', file);
        const xhr = new XMLHttpRequest();
        xhr.open('POST', `/api/upload?path=${encodeURIComponent(dir)}`);
        xhr.upload.onprogress = (e) => { if (e.lengthComputable && onProgress) onProgress(Math.round(e.loaded / e.total * 100)); };
        xhr.onload = () => xhr.status >= 200 && xhr.status < 300 ? resolve(xhr) : reject(new Error('Upload failed (' + xhr.status + ')'));
        xhr.onerror = () => reject(new Error('Upload failed'));
        xhr.send(f);
      });
    },
    space: () => mockEnabled ? mockSpace() : fetch('/api/space'),
    config: () => mockEnabled ? mockConfig() : fetch('/api/config'),
    delete: (path) => mockEnabled ? mockNotSupported() : fetch(`/api/files?path=${encodeURIComponent(path)}`, { method: 'DELETE' }),
    rename: (from, to) => mockEnabled ? mockNotSupported() : fetch('/api/rename', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ from, to }) }),
    mode: (forceRW) => mockEnabled ? mockNotSupported() : fetch('/api/mode', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ forceRW }) }),
  };

  const state = {
    path: '',
    entries: [],
    sort: { key: 'name', dir: 1 },
    loading: false,
    error: '',
    uploadProgress: -1, space: null,
    readonly: false,
    forceRW: false,
    pvcWatch: false,
    cluster: '',
    namespace: '',
    pvc: '',
    pod: '',
  };

  const PROJECT_INFO = {
    name: 'pvc-explorer-agent',
    description: 'Lightweight HTTP file-browser agent for browsing PersistentVolumeClaim contents through a simple REST API.',
    repoUrl: 'https://github.com/pvc-explorer-operator/pvc-explorer-agent',
    docsUrl: 'https://github.com/pvc-explorer-operator/pvc-explorer-agent#readme',
    license: 'Apache 2.0',
  };

  function setState(patch) { Object.assign(state, patch); render(); }

  function updateConfigState(cfg) {
    setState({
      readonly: !!cfg.readonly,
      forceRW: !!cfg.forceRW,
      pvcWatch: !!cfg.pvcWatch,
      cluster: cfg.cluster || '',
      namespace: cfg.namespace || '',
      pvc: cfg.pvc || '',
      pod: cfg.pod || '',
    });
  }

  function formatSize(bytes) {
    if (bytes == null) return '';
    if (bytes < 1024) return `${bytes} B`;
    const units = ['KB', 'MB', 'GB', 'TB'];
    let u = -1;
    do { bytes /= 1024; ++u; } while (bytes >= 1024 && u < units.length - 1);
    return `${bytes.toFixed(1)} ${units[u]}`;
  }

  function formatDate(ts) {
    if (!ts) return '';
    return new Date(ts).toLocaleString();
  }

  function join(base, name) {
    if (!base) return name;
    if (!name) return base;
    return base.replace(/\/+$/, '') + '/' + name.replace(/^\/+/, '');
  }

  function sortEntries(entries, { key, dir }) {
    const col = new Intl.Collator(undefined, { numeric: true, sensitivity: 'base' });
    return [...entries].sort((a, b) => {
      if (a.isDir !== b.isDir) return b.isDir - a.isDir;
      if (key === 'name') return col.compare(a.name, b.name) * dir;
      if (key === 'size') return ((a.size || 0) - (b.size || 0)) * dir;
      if (key === 'date') return (new Date(a.modTime) - new Date(b.modTime)) * dir;
      return 0;
    });
  }

  function el(tag, attrs, ...children) {
    const e = document.createElement(tag);
    if (attrs) Object.entries(attrs).forEach(([k, v]) => {
      if (k === 'cls') e.className = v;
      else if (k === 'text') e.textContent = v;
      else e[k] = v;
    });
    children.forEach(c => c && e.appendChild(c));
    return e;
  }

  function renderNavbar() {
    const nav = el('nav', { cls: 'navbar' });

    const brand = el('div', { cls: 'navbar-brand' });
    const logo = el('img', { cls: 'navbar-logo', src: '/logo.svg', alt: '' });
    logo.onerror = () => { logo.style.display = 'none'; };
    brand.appendChild(logo);
    const brandText = el('div', { cls: 'navbar-brand-text' });
    brandText.appendChild(el('div', { cls: 'navbar-brand-title', text: 'PVC Exporter Agent' }));
    brandText.appendChild(el('div', { cls: 'navbar-brand-sub', text: 'Storage Browser' }));
    brand.appendChild(brandText);
    nav.appendChild(brand);
    nav.appendChild(el('div', { cls: 'navbar-spacer' }));

    const pills = el('div', { cls: 'navbar-pills' });
    [
      state.cluster && { key: 'Cluster', value: state.cluster },
      state.namespace && { key: 'Namespace', value: state.namespace },
      state.pvc && { key: 'PVC', value: state.pvc },
      state.pod && { key: 'Pod', value: state.pod },
    ].filter(Boolean).forEach(f => {
      const pill = el('span', { cls: 'navbar-pill' });
      pill.appendChild(el('span', { cls: 'navbar-pill-key', text: f.key }));
      pill.appendChild(el('span', { cls: 'navbar-pill-value', text: f.value }));
      pills.appendChild(pill);
    });
    nav.appendChild(pills);

    if (state.pvcWatch) {
      if (state.readonly && !state.forceRW) {
        nav.appendChild(el('span', { cls: 'navbar-badge navbar-badge-warn', text: 'Read-Only' }));
      } else if (state.forceRW) {
        nav.appendChild(el('span', { cls: 'navbar-badge navbar-badge-danger', text: 'Forced RW' }));
      }
    }

    if (mockEnabled) {
      nav.appendChild(el('span', { cls: 'navbar-badge navbar-badge-warn', text: 'Mock Data' }));
    }

    const projectLinks = el('div', { cls: 'navbar-project-links' });
    const repoLink = el('a', {
      cls: 'navbar-project-link',
      href: PROJECT_INFO.repoUrl,
      target: '_blank',
      rel: 'noreferrer',
      title: PROJECT_INFO.name,
    });
    repoLink.setAttribute('aria-label', 'Open GitHub repository');
    repoLink.innerHTML = '<svg class="navbar-project-icon" viewBox="0 0 16 16" aria-hidden="true"><path fill="currentColor" d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.5-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82a7.65 7.65 0 0 1 4 0c1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.01 8.01 0 0 0 16 8c0-4.42-3.58-8-8-8Z"></path></svg>';
    projectLinks.appendChild(repoLink);
    nav.appendChild(projectLinks);

    return nav;
  }

  function renderTopbar() {
    const bar = el('div', { cls: 'topbar' });
    const parts = state.path ? state.path.split('/') : [];
    const rootCrumb = el('a', { cls: 'breadcrumb', tabIndex: 0, text: 'root' });
    rootCrumb.onclick = () => navigate('');
    rootCrumb.onkeydown = (e) => { if (e.key === 'Enter') navigate(''); };
    bar.appendChild(rootCrumb);
    let acc = '';
    parts.forEach((part, i) => {
      bar.appendChild(el('span', { cls: 'breadcrumb-sep', text: '/' }));
      acc = parts.slice(0, i + 1).join('/');
      const crumb = el('a', { cls: 'breadcrumb', tabIndex: 0, text: part });
      const path = acc;
      crumb.onclick = () => navigate(path);
      crumb.onkeydown = (e) => { if (e.key === 'Enter') navigate(path); };
      bar.appendChild(crumb);
    });
    return bar;
  }

  function renderNavBtns() {
    const wrap = el('div', { cls: 'nav-btns' });
    const btnUp = el('button', { type: 'button', cls: 'nav-btn', title: 'Up', tabIndex: 0, text: '↑' });
    const isRoot = !state.path || state.path === '/';
    if (isRoot) {
      btnUp.disabled = true;
    } else {
      btnUp.onclick = () => {
        const parts = state.path.split('/').filter(Boolean);
        parts.pop();
        navigate(parts.join('/'));
      };
    }
    wrap.appendChild(btnUp);
    return wrap;
  }

  function renderReadonlyBanner() {
    if (!state.pvcWatch) return el('div');
    if (state.readonly && !state.forceRW) {
      const banner = el('div', { cls: 'readonly-banner readonly-warning' });
      banner.appendChild(el('span', { text: 'Read-only — PVC is in use by another pod' }));
      const btn = el('button', { type: 'button', cls: 'readonly-btn override-btn', text: 'Override to Read-Write' });
      btn.disabled = !!state.loading;
      btn.onclick = () => {
        setState({ loading: true });
        API.mode(true)
          .then(r => r.ok ? API.config() : Promise.reject('Failed'))
          .then(r => r.json())
          .then(cfg => { updateConfigState(cfg); setState({ loading: false }); })
          .catch(e => setState({ loading: false, error: e.message || 'Failed to override' }));
      };
      banner.appendChild(btn);
      return banner;
    }
    if (!state.readonly && state.forceRW) {
      const banner = el('div', { cls: 'readonly-banner readonly-danger' });
      banner.appendChild(el('span', { text: 'Forced Read-Write active' }));
      const btn = el('button', { type: 'button', cls: 'readonly-btn revert-btn', text: 'Revert to Auto' });
      btn.disabled = !!state.loading;
      btn.onclick = () => {
        setState({ loading: true });
        API.mode(false)
          .then(r => r.ok ? API.config() : Promise.reject('Failed'))
          .then(r => r.json())
          .then(cfg => { updateConfigState(cfg); setState({ loading: false }); })
          .catch(e => setState({ loading: false, error: e.message || 'Failed to revert' }));
      };
      banner.appendChild(btn);
      return banner;
    }
    return el('div');
  }

  function renderDiskBar() {
    const div = el('div', { cls: 'header' });
    if (!state.space) { div.appendChild(renderSpinner()); return div; }
    const { used, total } = state.space;
    const percent = total ? Math.min(100, (used / total) * 100) : 0;
    const bar = el('div', { cls: 'disk-bar' });
    const fill = el('div', { cls: 'disk-bar-used' });
    fill.style.width = percent + '%';
    bar.appendChild(fill);
    bar.appendChild(el('div', { cls: 'disk-bar-label', text: `Disk: ${formatSize(used)} / ${formatSize(total)} (${percent.toFixed(1)}%)` }));
    div.appendChild(bar);
    return div;
  }

  function renderUploadBar() {
    if (state.readonly) return el('div');
    const div = el('div', { cls: 'upload-bar' });
    const input = el('input', { type: 'file', cls: 'upload-input', tabIndex: -1 });
    const btn = el('button', { cls: 'upload-btn', text: 'Upload File', tabIndex: 0 });

    btn.onclick = (e) => { e.preventDefault(); input.click(); };

    input.onchange = () => {
      if (!input.files || !input.files[0]) return;
      const file = input.files[0];
      input.value = '';
      setState({ uploadProgress: 0, error: '' });
      API.upload(state.path, file, (pct) => {
        setState({ uploadProgress: pct });
      })
        .then(() => { setState({ uploadProgress: -1 }); reload(); })
        .catch(e => setState({ uploadProgress: -1, error: e.message || 'Upload failed' }));
    };

    div.appendChild(input);
    div.appendChild(btn);

    if (state.uploadProgress >= 0) {
      const wrap = el('div', { cls: 'upload-progress-wrap' });
      const bar = el('div', { cls: 'upload-progress-bar' });
      bar.style.width = state.uploadProgress + '%';
      const pct = el('span', { cls: 'upload-progress-pct', text: state.uploadProgress + '%' });
      wrap.appendChild(bar);
      div.appendChild(wrap);
      div.appendChild(pct);
    }

    return div;
  }

  function renderFileTable() {
    const wrap = el('div', { cls: 'file-table-wrap' });
    wrap.appendChild(renderNavBtns());
    const table = el('table', { cls: 'file-table' });
    const thead = el('thead');
    const tr = el('tr');
    [{ key: 'name', label: 'Name' }, { key: 'size', label: 'Size' }, { key: 'date', label: 'Modified' }].forEach(col => {
      const th = el('th', { cls: 'sortable' + (state.sort.key === col.key ? ' sorted' : ''), tabIndex: 0 });
      th.textContent = col.label + (state.sort.key === col.key ? (state.sort.dir === 1 ? ' ▲' : ' ▼') : '');
      th.onclick = () => setSort(col.key);
      th.onkeydown = (e) => { if (e.key === 'Enter') setSort(col.key); };
      tr.appendChild(th);
    });
    tr.appendChild(el('th'));
    thead.appendChild(tr);
    table.appendChild(thead);

    const tbody = el('tbody');
    if (!state.entries.length) {
      const tdEmpty = el('td', { text: 'No files or folders.' });
      tdEmpty.colSpan = 4;
      tdEmpty.style.textAlign = 'center';
      tdEmpty.style.color = 'var(--text-muted)';
      tbody.appendChild(el('tr', {}, tdEmpty));
    } else {
      sortEntries(state.entries, state.sort).forEach(entry => {
        const row = el('tr', { cls: 'file-row' });
        const tdName = el('td', {}, el('span', { cls: 'file-icon', text: entry.isDir ? '📁' : '📄' }));
        const name = el('a', { cls: 'file-link', tabIndex: 0, text: entry.name });
        if (entry.isDir) {
          name.onclick = () => navigate(join(state.path, entry.name));
          name.onkeydown = (e) => { if (e.key === 'Enter') navigate(join(state.path, entry.name)); };
        } else {
          name.onclick = () => API.download(join(state.path, entry.name));
          name.onkeydown = (e) => { if (e.key === 'Enter') API.download(join(state.path, entry.name)); };
        }
        tdName.appendChild(name);
        const tdAct = el('td');
        if (!state.readonly) {
          const btnRename = el('button', { type: 'button', cls: 'action-btn rename-btn', title: 'Rename', tabIndex: 0 });
          btnRename.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>';
          btnRename.onclick = (e) => {
            e.stopPropagation();
            showModal({
              type: 'rename',
              initial: entry.name,
              onConfirm: (newName) => {
                if (newName && newName !== entry.name) {
                  setState({ loading: true, error: '' });
                  API.rename(join(state.path, entry.name), join(state.path, newName))
                    .then(r => { if (!r.ok) throw new Error('Rename failed'); reload(); })
                    .catch(e => setState({ loading: false, error: e.message || 'Rename failed' }));
                }
              }
            });
          };
          const btnDelete = el('button', { type: 'button', cls: 'action-btn delete-btn', title: 'Delete', tabIndex: 0 });
          btnDelete.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"/><path d="M10 11v6"/><path d="M14 11v6"/><path d="M9 6V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2"/></svg>';
          btnDelete.onclick = (e) => {
            e.stopPropagation();
            showModal({
              type: 'delete',
              message: `Delete ${entry.name}?`,
              onConfirm: () => {
                setState({ loading: true, error: '' });
                API.delete(join(state.path, entry.name))
                  .then(r => { if (!r.ok) throw new Error('Delete failed'); reload(); })
                  .catch(e => setState({ loading: false, error: e.message || 'Delete failed' }));
              }
            });
          };
          tdAct.appendChild(btnRename);
          tdAct.appendChild(btnDelete);
        }
        row.appendChild(tdName);
        row.appendChild(el('td', { text: entry.isDir ? '' : formatSize(entry.size) }));
        row.appendChild(el('td', { text: formatDate(entry.modTime) }));
        row.appendChild(tdAct);
        tbody.appendChild(row);
      });
    }
    table.appendChild(tbody);
    wrap.appendChild(table);
    return wrap;
  }

  function renderSpinner() {
    return el('div', { cls: 'spinner' }, el('div', { cls: 'spinner-anim' }));
  }

  function renderError(msg) {
    return el('div', { cls: 'error-bar' }, el('span', { text: '⚠' }), el('span', { text: msg }));
  }

  let navbarEl = null;
  let topbarEl = null;
  let mainEl = null;

  function ensureShell() {
    const root = document.getElementById('app-root');
    const boot = document.getElementById('boot-screen');
    if (boot) boot.remove();
    if (!navbarEl) {
      navbarEl = el('nav', { cls: 'navbar' });
      document.body.insertBefore(navbarEl, document.body.firstChild);
    }
    if (!topbarEl) {
      topbarEl = el('div', { cls: 'topbar' });
      root.appendChild(topbarEl);
    }
    if (!mainEl) {
      mainEl = el('main', { cls: 'main-content' });
      root.appendChild(mainEl);
    }
  }

  function render() {
    ensureShell();

    const newNav = renderNavbar();
    navbarEl.replaceWith(newNav);
    navbarEl = newNav;

    const newTopbar = renderTopbar();
    topbarEl.replaceWith(newTopbar);
    topbarEl = newTopbar;

    mainEl.innerHTML = '';
    mainEl.appendChild(renderReadonlyBanner());
    mainEl.appendChild(renderDiskBar());
    if (state.error) mainEl.appendChild(renderError(state.error));
    mainEl.appendChild(renderUploadBar());
    mainEl.appendChild(state.loading ? renderSpinner() : renderFileTable());
  }

  let retryTimer = null;

  function showModal(opts) {
    const overlay = el('div', { cls: 'modal-overlay' });
    const modal = el('div', { cls: 'modal-box' });
    let input = null;
    let confirmed = false;
    if (opts.type === 'rename') {
      modal.appendChild(el('div', { cls: 'modal-title', text: 'Rename' }));
      input = el('input', {
        cls: 'modal-input',
        type: 'text',
        value: opts.initial || '',
        autofocus: true,
        spellcheck: false,
        autocomplete: 'off',
        maxlength: 255
      });
      modal.appendChild(input);
      setTimeout(() => input.focus(), 10);
    }
    if (opts.type === 'delete') {
      modal.appendChild(el('div', { cls: 'modal-title', text: opts.message || 'Are you sure?' }));
    }
    const btnRow = el('div', { cls: 'modal-btn-row' });
    const btnConfirm = el('button', { type: 'button', cls: 'modal-btn modal-btn-confirm', text: opts.type === 'rename' ? 'Rename' : 'Delete' });
    const btnCancel = el('button', { type: 'button', cls: 'modal-btn modal-btn-cancel', text: 'Cancel' });
    btnConfirm.onclick = () => {
      confirmed = true;
      close();
      if (opts.onConfirm) opts.onConfirm(input ? input.value : undefined);
    };
    btnCancel.onclick = () => {
      close();
      if (opts.onCancel) opts.onCancel();
    };
    btnRow.appendChild(btnConfirm);
    btnRow.appendChild(btnCancel);
    modal.appendChild(btnRow);
    function onKey(e) {
      if (e.key === 'Escape') {
        close();
        if (opts.onCancel) opts.onCancel();
      }
      if (e.key === 'Enter' && (opts.type === 'delete' || document.activeElement === input)) {
        confirmed = true;
        close();
        if (opts.onConfirm) opts.onConfirm(input ? input.value : undefined);
      }
    }
    function close() {
      document.removeEventListener('keydown', onKey, true);
      overlay.remove();
    }
    document.addEventListener('keydown', onKey, true);
    overlay.appendChild(modal);
    document.body.appendChild(overlay);
  }

  function renderWaiting() {
    ensureShell();
    const newNav = renderNavbar();
    navbarEl.replaceWith(newNav);
    navbarEl = newNav;
    topbarEl.innerHTML = '';
    mainEl.innerHTML = '';
    mainEl.appendChild(el('div', { cls: 'waiting-screen' },
      el('div', { cls: 'waiting-icon', text: '⏳' }),
      el('div', { cls: 'waiting-title', text: 'PVC is in use' }),
      el('div', { cls: 'waiting-msg', text: 'The volume is currently attached to another pod (cronjob running). Waiting for it to finish…' }),
      el('div', { cls: 'waiting-sub', text: 'This page will refresh automatically.' })
    ));
  }

  function navigate(path) {
    setState({ path, loading: true, error: '' });
    Promise.all([
      API.list(path).then(r => { if (!r.ok) throw new Error('Failed to list'); return r.json(); }),
      state.space ? Promise.resolve(state.space) : API.space().then(r => { if (!r.ok) throw new Error('Failed to get disk usage'); return r.json(); }),
    ])
      .then(([files, space]) => setState({ entries: files.entries || [], loading: false, error: '', space }))
      .catch(e => setState({ loading: false, error: e.message || 'Failed to load' }));
  }

  function reload() { navigate(state.path); }

  function setSort(key) {
    setState({ sort: { key, dir: state.sort.key === key ? -state.sort.dir : 1 } });
  }

  function reloadConfigAndMaybeNavigate() {
    setState({ loading: true });
    API.config()
      .then(r => r.ok ? r.json() : {})
      .then(cfg => {
        clearTimeout(retryTimer);
        updateConfigState(cfg);
        navigate(state.path || '');
      })
      .catch(() => {
        if (!mockEnabled && allowAutoMock) {
          mockEnabled = true;
          reloadConfigAndMaybeNavigate();
          return;
        }
        renderWaiting();
        retryTimer = setTimeout(reloadConfigAndMaybeNavigate, 10000);
      });
  }

  reloadConfigAndMaybeNavigate();

  document.getElementById('app-root').addEventListener('keydown', (e) => {
    if (e.key !== 'ArrowDown' && e.key !== 'ArrowUp') return;
    const links = document.querySelectorAll('.file-link');
    if (!links.length) return;
    let idx = Array.prototype.indexOf.call(links, document.activeElement);
    idx = e.key === 'ArrowDown' ? Math.min(links.length - 1, idx + 1) : Math.max(0, idx - 1);
    links[idx].focus();
    e.preventDefault();
  });
})();
