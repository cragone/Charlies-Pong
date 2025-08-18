import React, { createContext } from "react";
import type {ReactNode} from "react";

interface GlobalContextType {
    backendPath: string;
}


const GlobalContext = createContext<GlobalContextType | undefined>(undefined);


interface GlobalProviderProps {
    children: ReactNode;
}

const GlobalProvider: React.FC<GlobalProviderProps>= ({children}) => {
    const backendPath = import.meta.env.VITE_API_BASE;

    const contextValue: GlobalContextType = {
        backendPath,
    }

    return (
        <GlobalContext.Provider value={contextValue}>
            {children}
        </GlobalContext.Provider>
    )
}
export default GlobalContext;

export {GlobalContext, GlobalProvider}