import * as React from "react";
import { createRoot } from "react-dom/client";
import App from "./components/App";
import { Provider } from "react-redux";
import store from "./redux/store";


const root = createRoot(document.getElementById("root"));
root.render(<Provider store={store}><App /></Provider>);
