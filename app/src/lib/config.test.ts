import { afterEach, describe, expect, it, vi } from 'vitest';

const VALID_ENV = {
  VITE_API_URL: 'https://api.dagger-ai.com',
  VITE_SUPABASE_URL: 'https://example.supabase.co',
  VITE_SUPABASE_ANON_KEY: 'anon-key',
};

async function importConfigModule(env: Record<string, string> = VALID_ENV) {
  for (const [key, value] of Object.entries(env)) {
    vi.stubEnv(key, value);
  }
  vi.resetModules();
  return import('#/lib/config');
}

afterEach(() => {
  vi.unstubAllEnvs();
  vi.resetModules();
});

describe('loadConfig', () => {
  it('returns the three values when all are present', () => {
    return importConfigModule().then(({ loadConfig }) => {
      expect(loadConfig(VALID_ENV)).toEqual({
        apiUrl: 'https://api.dagger-ai.com',
        supabaseUrl: 'https://example.supabase.co',
        supabaseAnonKey: 'anon-key',
      });
    });
  });

  it('names the missing variable', async () => {
    const { loadConfig } = await importConfigModule();

    const { VITE_API_URL: _omitted, ...incomplete } = VALID_ENV;

    expect(() => loadConfig(incomplete)).toThrow(/VITE_API_URL/);
  });

  it('names every missing variable in a single error', async () => {
    const { loadConfig } = await importConfigModule();

    let thrown: unknown;
    try {
      loadConfig({ VITE_API_URL: VALID_ENV.VITE_API_URL });
    } catch (error) {
      thrown = error;
    }

    expect(thrown).toBeInstanceOf(Error);
    expect((thrown as Error).message).toContain('VITE_SUPABASE_URL');
    expect((thrown as Error).message).toContain('VITE_SUPABASE_ANON_KEY');
  });

  it('treats an empty string as missing', async () => {
    const { loadConfig } = await importConfigModule();

    expect(() => loadConfig({ ...VALID_ENV, VITE_API_URL: '' })).toThrow(
      /VITE_API_URL/,
    );
  });
});

describe('config', () => {
  it('is read from the environment once, at import time', async () => {
    const { config } = await importConfigModule();

    expect(config.apiUrl).toBe('https://api.dagger-ai.com');
    expect(config.supabaseUrl).toBe('https://example.supabase.co');
    expect(config.supabaseAnonKey).toBe('anon-key');
  });

  it('fails at import time when a variable is missing', async () => {
    await expect(
      importConfigModule({
        ...VALID_ENV,
        VITE_SUPABASE_ANON_KEY: '',
      }),
    ).rejects.toThrow(/VITE_SUPABASE_ANON_KEY/);
  });
});
