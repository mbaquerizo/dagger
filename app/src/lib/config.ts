import z from 'zod';

const envSchema = z
  .object({
    VITE_API_URL: z.string().trim().min(1),
    VITE_SUPABASE_URL: z.string().trim().min(1),
    VITE_SUPABASE_ANON_KEY: z.string().trim().min(1),
  })
  .transform((arg) => ({
    apiUrl: arg.VITE_API_URL,
    supabaseUrl: arg.VITE_SUPABASE_URL,
    supabaseAnonKey: arg.VITE_SUPABASE_ANON_KEY,
  }));

export type Config = z.infer<typeof envSchema>;

export const loadConfig = (env: Record<string, string | undefined>) => {
  const result = envSchema.safeParse(env);

  if (!result.success) {
    const errorMessage = result.error.issues
      .map((issue) => issue.path.join('.'))
      .join(', ');

    throw new Error(errorMessage);
  }

  return result.data;
};

export const config: Config = loadConfig(import.meta.env);
