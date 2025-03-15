import { useState, useEffect } from "react";
import confetti from "canvas-confetti";
import { MapPin, Frown } from "lucide-react";
import { Question, AnswerResponse } from "../types";
import toast from "react-hot-toast";

interface GameProps {
    accessToken: string;
    onScoreUpdate: () => void;
}

export function Game({ accessToken, onScoreUpdate }: GameProps) {
    const [questions, setQuestions] = useState < Question[] > ([]);
    const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
    const [answerResult, setAnswerResult] = useState < AnswerResponse | null > (null);
    const [loading, setLoading] = useState(false);

    const currentQuestion = questions[currentQuestionIndex];

    useEffect(() => {
        fetchQuestions();
    }, []);

    const fetchQuestions = async () => {
        try {
            const response = await fetch(
                `${import.meta.env.VITE_BACKEND_API}/questions`,
                {
                    headers: {
                        Authorization: `Bearer ${accessToken}`,
                    },
                }
            );
            if (!response.ok) throw new Error("Failed to fetch questions");
            const data = await response.json();
            setQuestions(data.data);
        } catch (error) {
            toast.error("Failed to load questions");
        }
    };

    const handleAnswer = async (answer: string) => {
        if (loading || !currentQuestion) return;
        setLoading(true);
        try {
            const response = await fetch(
                `${import.meta.env.VITE_BACKEND_API}/questions/check`,
                {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        Authorization: `Bearer ${accessToken}`,
                    },
                    body: JSON.stringify({
                        question_id: currentQuestion.question_id,
                        answer,
                    }),
                }
            );
            if (!response.ok) throw new Error("Failed to check answer");
            const result = await response.json();
            setAnswerResult(result.data);
            onScoreUpdate();

            if (result.data.correct) {
                confetti({
                    particleCount: 100,
                    spread: 70,
                    origin: { y: 0.6 },
                });
            }
        } catch (error) {
            console.error("Error submitting answer:", error);
            toast.error("Failed to submit answer");
        } finally {
            setLoading(false);
        }
    };

    const handleNext = () => {
        if (currentQuestionIndex === questions.length - 1) {
            fetchQuestions();
            setCurrentQuestionIndex(0);
        } else {
            setCurrentQuestionIndex((prev) => prev + 1);
        }
        setAnswerResult(null);
    };

    if (!currentQuestion) {
        return (
            <div className="flex items-center justify-center h-96">
                <div className="animate-spin rounded-full h-16 w-16 border-t-2 border-b-2 border-blue-500"></div>
            </div>
        );
    }

    return (
        <div className="max-w-4xl mx-auto p-6">
            <div className="bg-white rounded-xl shadow-lg p-8">
                {!answerResult ? (
                    <>
                        <div className="space-y-6 mb-8">
                            <div className="flex items-start space-x-4">
                                <MapPin className="w-6 h-6 text-blue-500 flex-shrink-0 mt-1" />
                                <div>
                                    <h2 className="text-2xl font-bold mb-4">Where am I?</h2>
                                    <ul className="space-y-3">
                                        {currentQuestion.question_hints.map((hint, index) => (
                                            <li
                                                key={index}
                                                className="text-gray-700 bg-blue-50 p-3 rounded-lg"
                                            >
                                                {hint}
                                            </li>
                                        ))}
                                    </ul>
                                </div>
                            </div>
                        </div>

                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            {currentQuestion.answer_options.map((option, index) => (
                                <button
                                    key={index}
                                    onClick={() => handleAnswer(option)}
                                    disabled={loading}
                                    className="p-4 text-left rounded-lg border-2 border-gray-200 hover:border-blue-500 hover:bg-blue-50 transition-colors disabled:opacity-50"
                                >
                                    {option}
                                </button>
                            ))}
                        </div>
                    </>
                ) : (
                    <div className="space-y-6">
                        <div className="flex items-center justify-center">
                            {answerResult.correct ? (
                                <div className="text-center">
                                    <h3 className="text-2xl font-bold text-green-600 mb-2">
                                        Correct! 🎉
                                    </h3>
                                    <p className="text-gray-600">
                                        You've found {answerResult.correct_answer}!
                                    </p>
                                </div>
                            ) : (
                                <div className="text-center">
                                    <Frown className="w-16 h-16 text-red-500 mx-auto mb-2" />
                                    <h3 className="text-2xl font-bold text-red-600 mb-2">
                                        Not quite!
                                    </h3>
                                    <p className="text-gray-600">
                                        The correct answer was {answerResult.correct_answer}
                                    </p>
                                </div>
                            )}
                        </div>

                        <div className="space-y-4">
                            <div>
                                <h4 className="font-bold text-lg mb-2">Fun Facts</h4>
                                <ul className="space-y-2">
                                    {answerResult.fun_facts.map((fact, index) => (
                                        <li
                                            key={index}
                                            className="bg-purple-50 p-3 rounded-lg text-gray-700"
                                        >
                                            {fact}
                                        </li>
                                    ))}
                                </ul>
                            </div>

                            <div>
                                <h4 className="font-bold text-lg mb-2">Trivia</h4>
                                <ul className="space-y-2">
                                    {answerResult.trivia.map((item, index) => (
                                        <li
                                            key={index}
                                            className="bg-blue-50 p-3 rounded-lg text-gray-700"
                                        >
                                            {item}
                                        </li>
                                    ))}
                                </ul>
                            </div>
                        </div>

                        <button
                            onClick={handleNext}
                            className="w-full bg-gradient-to-r from-blue-500 to-purple-600 text-white py-3 rounded-lg font-medium hover:opacity-90 transition-opacity"
                        >
                            Next Destination
                        </button>
                    </div>
                )}
            </div>
        </div>
    );
}
