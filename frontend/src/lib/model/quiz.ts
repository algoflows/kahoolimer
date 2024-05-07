export interface Quiz {
  id: string;
  name: string;
  questions: QuizQuestion[];
}

export interface Player {
  id: string;
  name: string;
}

export interface QuizQuestion {
  id: string;
  question: string;
  options: string[];
}

export interface QuizChoice {
  id: string;
  name: string;
  correct: boolean;
}
