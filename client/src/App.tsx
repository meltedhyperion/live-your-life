import { useState, useEffect } from "react";
import { Toaster } from "react-hot-toast";
import toast from "react-hot-toast";
import { supabase } from "./lib/supabase";
import { AuthForm } from "./components/AuthForm";
import { Navbar } from "./components/Navbar";
import { Game } from "./components/Game";
import { Leaderboard } from "./components/Leaderboard";
import { Player, LeaderboardResponse } from "./types";

const useImageSlider = (interval = 5000) => {
  const images = [
    "scenes/1.png",
    "scenes/2.png",
    "scenes/3.png",
    "scenes/4.png",
    "scenes/5.png",
    "scenes/6.png",
    "scenes/7.png",
    "scenes/8.png",
  ];
  const [currentImage, setCurrentImage] = useState(images[0]);

  useEffect(() => {
    const imageInterval = setInterval(() => {
      setCurrentImage((prevImage) => {
        const currentIndex = images.indexOf(prevImage);
        const nextIndex = (currentIndex + 1) % images.length;
        return images[nextIndex];
      });
    }, interval);

    return () => clearInterval(imageInterval);
  }, [images, interval]);

  return currentImage;
};

function App() {
  const [session, setSession] = useState<any>(null);
  const [player, setPlayer] = useState<Player | null>(null);
  const [showLeaderboard, setShowLeaderboard] = useState(false);
  const [leaderboard, setLeaderboard] = useState<LeaderboardResponse | null>(
    null
  );
  const [isLoading, setIsLoading] = useState(true);
  const currentBackground = useImageSlider();

  // Fetch session and player on mount
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
          // Check if player profile exists
          const response = await fetch(
            `${import.meta.env.VITE_BACKEND_API}/players`,
            {
              headers: {
                Authorization: `Bearer ${session.access_token}`,
              },
            }
          );

          if (!response.ok && response.status === 404) {
            // If player profile doesn't exist, create one
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
                }),
              }
            );

            if (!createResponse.ok) {
              throw new Error("Failed to create player profile");
            }
          }

          await fetchPlayer(session.access_token);
        } catch (error) {
          // Uncomment to debug errors
          // console.error("Error handling player profile:", error);
          // toast.error("Failed to set up player profile");
        }
      }
      setIsLoading(false);
    });

    return () => subscription.unsubscribe();
  }, []);

  // Handle invite code (runs whenever session changes)
  useEffect(() => {
    const handleInviteCode = async () => {
      const urlParams = new URLSearchParams(window.location.search);
      const inviteCode = urlParams.get("invite-code");

      if (inviteCode && session) {
        try {
          const response = await fetch(
            `${import.meta.env.VITE_BACKEND_API}/friends/${inviteCode}`,
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
          // Remove the invite code from the URL
          window.history.replaceState({}, "", window.location.pathname);
        } catch (error) {
          // console.error("Error adding friend:", error);
          // toast.error("Failed to add friend");
        }
      }
    };

    handleInviteCode();
  }, [session]);

  // Function to fetch player data using access token
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
      if (!player) {
        // toast.error("Failed to load player profile");
      }
    }
  };

  const fetchLeaderboard = async () => {
    try {
      const response = await fetch(
        `${import.meta.env.VITE_BACKEND_API}/players/leaderboard`,
        {
          headers: {
            Authorization: `Bearer ${session?.access_token}`,
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

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (!session) {
    return (
      <>
        <AuthForm onAuthSuccess={() => {}} />
        <Toaster position="top-center" />
      </>
    );
  }

  return (
    <div
      className="min-h-screen bg-cover bg-no-repeat bg-center transition-opacity duration-1000 ease-in-out"
      style={{ backgroundImage: `url(${currentBackground})` }}
    >
      <div className="min-h-screen">
        <Toaster position="top-center" />
        {/* Always render Navbar once logged in; Navbar can handle a null player gracefully */}
        <Navbar
          player={player}
          onShowLeaderboard={fetchLeaderboard}
          onLogout={handleLogout}
        />
        <main className="py-8">
          <Game
            accessToken={session.access_token}
            onScoreUpdate={() => fetchPlayer(session.access_token)}
          />
        </main>
        {showLeaderboard && leaderboard && (
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
