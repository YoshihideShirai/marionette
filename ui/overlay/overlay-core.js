const DEFAULT_PORTAL_ID = 'ui-overlay-root';

const state = {
  openStack: [],
  scrollLockCount: 0,
  savedScrollY: 0,
};

function getFocusable(container) {
  const selector = [
    'a[href]',
    'button:not([disabled])',
    'input:not([disabled]):not([type="hidden"])',
    'select:not([disabled])',
    'textarea:not([disabled])',
    '[tabindex]:not([tabindex="-1"])',
    '[contenteditable="true"]',
  ].join(',');
  return [...container.querySelectorAll(selector)].filter((el) => !el.hasAttribute('inert'));
}

function resolveInitialFocus(panel, initialFocus) {
  if (!initialFocus) return getFocusable(panel)[0] || panel;
  if (typeof initialFocus === 'string') return panel.querySelector(initialFocus) || panel;
  return initialFocus;
}

function lockScroll() {
  if (state.scrollLockCount > 0) {
    state.scrollLockCount += 1;
    return;
  }
  state.scrollLockCount = 1;
  state.savedScrollY = window.scrollY;
  document.body.dataset.overlayScrollLock = 'true';
  document.body.style.position = 'fixed';
  document.body.style.top = `-${state.savedScrollY}px`;
  document.body.style.left = '0';
  document.body.style.right = '0';
  document.body.style.width = '100%';
  document.body.style.overflow = 'hidden';
}

function unlockScroll() {
  if (state.scrollLockCount === 0) return;
  state.scrollLockCount -= 1;
  if (state.scrollLockCount > 0) return;

  document.body.removeAttribute('data-overlay-scroll-lock');
  document.body.style.position = '';
  document.body.style.top = '';
  document.body.style.left = '';
  document.body.style.right = '';
  document.body.style.width = '';
  document.body.style.overflow = '';
  window.scrollTo(0, state.savedScrollY);
}

export function ensureOverlayPortalRoot(id = DEFAULT_PORTAL_ID) {
  let portal = document.getElementById(id);
  if (portal) return portal;

  portal = document.createElement('div');
  portal.id = id;
  portal.className = 'ui-overlay-portal-root';
  portal.setAttribute('aria-live', 'polite');
  document.body.appendChild(portal);
  return portal;
}

function createLayer({
  type,
  panel,
  closeOnEsc,
  closeOnBackdrop,
  lockBodyScroll,
  zIndex,
  onRequestClose,
  restoreFocusTo,
  initialFocus,
  portal,
}) {
  const layer = document.createElement('div');
  layer.className = `ui-overlay-layer ui-overlay-layer-${type}`;
  layer.style.zIndex = String(zIndex);

  const backdrop = document.createElement('button');
  backdrop.type = 'button';
  backdrop.className = 'ui-overlay-backdrop';
  backdrop.setAttribute('aria-label', `${type} backdrop`);

  panel.classList.add('ui-overlay-panel');
  panel.tabIndex = panel.tabIndex >= 0 ? panel.tabIndex : -1;

  layer.appendChild(backdrop);
  layer.appendChild(panel);
  portal.appendChild(layer);

  const close = (reason = 'programmatic') => {
    const idx = state.openStack.findIndex((entry) => entry.layer === layer);
    if (idx >= 0) state.openStack.splice(idx, 1);
    layer.remove();
    if (lockBodyScroll) unlockScroll();
    if (restoreFocusTo && typeof restoreFocusTo.focus === 'function') {
      restoreFocusTo.focus({ preventScroll: true });
    }
    onRequestClose?.(reason);
  };

  const onBackdropClick = () => {
    if (closeOnBackdrop) close('backdrop');
  };

  const onKeyDown = (event) => {
    const top = state.openStack[state.openStack.length - 1];
    if (!top || top.layer !== layer) return;

    if (event.key === 'Escape' && closeOnEsc) {
      event.preventDefault();
      close('escape');
      return;
    }

    if (event.key !== 'Tab') return;
    const focusable = getFocusable(panel);
    if (focusable.length === 0) {
      event.preventDefault();
      panel.focus({ preventScroll: true });
      return;
    }

    const first = focusable[0];
    const last = focusable[focusable.length - 1];
    const active = document.activeElement;
    if (event.shiftKey && active === first) {
      event.preventDefault();
      last.focus({ preventScroll: true });
    } else if (!event.shiftKey && active === last) {
      event.preventDefault();
      first.focus({ preventScroll: true });
    }
  };

  backdrop.addEventListener('click', onBackdropClick);
  layer.addEventListener('keydown', onKeyDown);

  state.openStack.push({ layer, close });

  if (lockBodyScroll) lockScroll();

  queueMicrotask(() => {
    const target = resolveInitialFocus(panel, initialFocus);
    target.focus({ preventScroll: true });
  });

  return { close, layer, panel };
}

export function createOverlayCore(config = {}) {
  const portal = ensureOverlayPortalRoot(config.portalId || DEFAULT_PORTAL_ID);

  return {
    open({
      type = 'modal',
      panel,
      closeOnEsc = true,
      closeOnBackdrop = true,
      lockBodyScroll = type === 'modal' || type === 'drawer',
      zIndex = 1000,
      onRequestClose,
      restoreFocusTo = document.activeElement,
      initialFocus,
    }) {
      return createLayer({
        type,
        panel,
        closeOnEsc,
        closeOnBackdrop,
        lockBodyScroll,
        zIndex,
        onRequestClose,
        restoreFocusTo,
        initialFocus,
        portal,
      });
    },
    closeTop(reason = 'programmatic') {
      const top = state.openStack[state.openStack.length - 1];
      if (!top) return;
      top.close(reason);
    },
    closeAll(reason = 'programmatic') {
      while (state.openStack.length > 0) {
        const top = state.openStack[state.openStack.length - 1];
        top.close(reason);
      }
    },
  };
}
