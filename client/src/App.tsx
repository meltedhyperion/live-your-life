import { useState, useEffect } from "react";
import { Toaster } from "react-hot-toast";
import toast from "react-hot-toast";
import { supabase } from "./lib/supabase";
import { AuthForm } from "./components/AuthForm";
import { Navbar } from "./components/Navbar";
import { Game } from "./components/Game";
import { Leaderboard } from "./components/Leaderboard";
import { Player, LeaderboardResponse } from "./types";
import { useImageSlider } from "./lib/utils";

function App() {
  const [session, setSession] = useState<any>(null);
  const [player, setPlayer] = useState<Player | null>(null);
  const [showLeaderboard, setShowLeaderboard] = useState(false);
  const [leaderboard, setLeaderboard] = useState<LeaderboardResponse | null>(
    null
  );
  const [isLoading, setIsLoading] = useState(true);
  const [guestStats, setGuestStats] = useState({ correct: 0, total: 0 });
  const [showAuthForm, setShowAuthForm] = useState(false);
  const currentBackground = useImageSlider();

  useEffect(() => {
    const storedTotal = sessionStorage.getItem("guest_total");
    const storedCorrect = sessionStorage.getItem("guest_correct");
    if (storedTotal && storedCorrect) {
      setGuestStats({
        total: Number(storedTotal),
        correct: Number(storedCorrect),
      });
    }
  }, []);

  useEffect(() => {
    const fetchUserAndPlayer = async () => {
      const {
        data: { session },
      } = await supabase.auth.getSession();
      setSession(session);
      if (session) {
        await fetchPlayer(session.access_token);
      }
      setIsLoading(false);
    };

    fetchUserAndPlayer();

    const {
      data: { subscription },
    } = supabase.auth.onAuthStateChange(async (_event, session) => {
      setSession(session);
      if (session) {
        try {
          const response = await fetch(
            `${import.meta.env.VITE_BACKEND_API}/players`,
            {
              headers: {
                Authorization: `Bearer ${session.access_token}`,
              },
            }
          );

          if (!response.ok && response.status === 404) {
            const createResponse = await fetch(
              `${import.meta.env.VITE_BACKEND_API}/players/create`,
              {
                method: "POST",
                headers: {
                  "Content-Type": "application/json",
                  Authorization: `Bearer ${session.access_token}`,
                },
                body: JSON.stringify({
                  name: session.user.email?.split("@")[0],
                  total: Number(sessionStorage.getItem("guest_total") || "0"),
                  correct: Number(
                    sessionStorage.getItem("guest_correct") || "0"
                  ),
                }),
              }
            );

            if (!createResponse.ok) {
              throw new Error("Failed to create player profile");
            }
          }

          await fetchPlayer(session.access_token);
        } catch (error) {
          // console.error("Error handling player profile:", error);
          // toast.error("Failed to set up player profile");
        }
      }
      setIsLoading(false);
    });

    return () => subscription.unsubscribe();
  }, []);

  useEffect(() => {
    const handleInviteCode = async () => {
      const urlParams = new URLSearchParams(window.location.search);
      const inviteCodeFromUrl = urlParams.get("invite-code");

      if (inviteCodeFromUrl) {
        localStorage.setItem("inviteCode", inviteCodeFromUrl);
        window.history.replaceState({}, "", window.location.pathname);
      }

      const storedInviteCode = localStorage.getItem("inviteCode");
      if (storedInviteCode && session) {
        try {
          const response = await fetch(
            `${import.meta.env.VITE_BACKEND_API}/friends/${storedInviteCode}`,
            {
              method: "POST",
              headers: {
                Authorization: `Bearer ${session.access_token}`,
              },
            }
          );

          if (!response.ok) {
            throw new Error("Failed to add friend");
          }

          toast.success("Friend added successfully!");
          localStorage.removeItem("inviteCode");
        } catch (error) {
          // console.error("Error adding friend:", error);
          // toast.error("Failed to add friend");
          localStorage.removeItem("inviteCode");
        }
      }
    };

    handleInviteCode();
  }, [session]);

  const fetchPlayer = async (accessToken: string) => {
    try {
      const response = await fetch(
        `${import.meta.env.VITE_BACKEND_API}/players`,
        {
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        }
      );

      if (!response.ok) throw new Error("Failed to fetch player data");

      const data = await response.json();
      setPlayer(data.data);
    } catch (error) {
      console.error("Error fetching player:", error);
    }
  };

  const fetchLeaderboard = async () => {
    if (!session) return;
    try {
      const response = await fetch(
        `${import.meta.env.VITE_BACKEND_API}/players/leaderboard`,
        {
          headers: {
            Authorization: `Bearer ${session.access_token}`,
          },
        }
      );

      const data = await response.json();

      if (response.status === 404) {
        toast.error(data.message || "No friends found in leaderboard");
        return;
      }

      if (!response.ok) {
        throw new Error("Failed to fetch leaderboard");
      }

      setLeaderboard(data.data);
      setShowLeaderboard(true);
    } catch (error) {
      console.error("Error fetching leaderboard:", error);
      toast.error("Failed to load leaderboard");
    }
  };

  const handleLogout = async () => {
    await supabase.auth.signOut();
    setSession(null);
    setPlayer(null);
  };

  const updateScore = () => {
    if (session) {
      fetchPlayer(session.access_token);
    } else {
      const total = Number(sessionStorage.getItem("guest_total") || "0");
      const correct = Number(sessionStorage.getItem("guest_correct") || "0");
      setGuestStats({ total, correct });
    }
  };

  if (!session && localStorage.getItem("inviteCode")) {
    return (
      <>
        <AuthForm onAuthSuccess={() => setShowAuthForm(false)} />
        <Toaster position="top-center" />
      </>
    );
  }

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (showAuthForm) {
    return <AuthForm onAuthSuccess={() => setShowAuthForm(false)} />;
  }

  return (
    <div
      className="min-h-screen bg-cover bg-no-repeat bg-center transition-opacity duration-1000 ease-in-out"
      style={{ backgroundImage: `url(${currentBackground})` }}
    >
      <div className="min-h-screen">
        <Toaster position="top-center" />
        <Navbar
          player={player}
          isGuest={!session}
          guestStats={guestStats}
          onShowAuthForm={() => setShowAuthForm(true)}
          onShowLeaderboard={fetchLeaderboard}
          onLogout={handleLogout}
        />
        <main className="py-8">
          <Game
            accessToken={session ? session.access_token : null}
            onScoreUpdate={updateScore}
          />
        </main>
        {session && showLeaderboard && leaderboard && (
          <Leaderboard
            leaderboard={leaderboard.player_stats}
            onClose={() => setShowLeaderboard(false)}
          />
        )}
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

export default App;
