export interface Player {
  id: string;
  name: string;
  avatar: string;
  correct_answers: number;
  total_attempts: number;
  score: number;
  created_at: string;
  updated_at: string;
}

export interface Question {
  question_id: number;
  question_hints: string[];
  answer_options: string[];
}

export interface AnswerResponse {
  correct: boolean;
  fun_facts: string[];
  trivia: string[];
  correct_answer: string;
  correct_answers: number;
  total_attempts: number;
  score: number;
}

export interface LeaderboardEntry {
  name: string;
  avatar: string;
  score: number;
  correct_answers: number;
  total_attempts: number;
}

export interface LeaderboardResponse {
  player_stats: LeaderboardEntry[];
}