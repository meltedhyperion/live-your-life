import { useState } from "react";
import { supabase } from "../lib/supabase";
import { Globe2 } from "lucide-react";
import toast from "react-hot-toast";

interface AuthFormProps {
  onAuthSuccess: () => void;
}

export function AuthForm({ onAuthSuccess }: AuthFormProps) {
  const [isLogin, setIsLogin] = useState(true);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      if (isLogin) {
        const { error } = await supabase.auth.signInWithPassword({
          email,
          password,
        });
        if (error) throw error;
        toast.success("Welcome back!");
        onAuthSuccess();
      } else {
        // Sign up flow
        const { data: signUpData, error: signUpError } =
          await supabase.auth.signUp({
            email,
            password,
            options: {
              emailRedirectTo: `${window.location.origin}`,
            },
          });

        if (signUpError) throw signUpError;

        if (!signUpData.session) {
          // Email confirmation is required
          toast.success("Please check your email to confirm your account!");
        }
      }
    } catch (error: any) {
      toast.error(error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center p-4">
      <div className="bg-white rounded-xl shadow-2xl p-8 w-full max-w-md">
        <div className="flex items-center justify-center mb-8">
          <Globe2 className="w-12 h-12 text-blue-500" />
          <h1 className="text-3xl font-bold ml-3 bg-gradient-to-r from-blue-500 to-purple-600 bg-clip-text text-transparent">
            Globetrotter
          </h1>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Email
            </label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Password
            </label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              required
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-gradient-to-r from-blue-500 to-purple-600 text-white py-2 rounded-lg font-medium hover:opacity-90 transition-opacity disabled:opacity-50"
          >
            {loading ? "Loading..." : isLogin ? "Sign In" : "Sign Up"}
          </button>
        </form>

        <p className="mt-4 text-center text-sm text-gray-600">
          {isLogin ? "Don't have an account? " : "Already have an account? "}
          <button
            onClick={() => setIsLogin(!isLogin)}
            className="text-blue-500 hover:underline font-medium"
          >
            {isLogin ? "Sign Up" : "Sign In"}
          </button>
        </p>
      </div>
      <footer className="fixed bottom-0 left-0 right-0 text-center p-4 text-sm text-white">
        Made with{" "}
        <span role="img" aria-label="heart" className="text-red-500">
          ❤️
        </span>{" "}
        by aryansingh.dev
      </footer>
    </div>
  );
}
