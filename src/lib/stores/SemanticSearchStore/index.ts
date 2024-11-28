import { create } from "zustand";

interface State {

}

const initialState: State = {

}

interface SemanticSearchStoreState extends State {

}

export const useSemanticSearchStore = create<SemanticSearchStoreState>((set, get) => ({
  ...initialState,

}))
