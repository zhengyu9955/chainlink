import { createMuiTheme, MuiThemeProvider } from '@material-ui/core/styles'
import JavascriptTimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'
import moment from 'moment'
import React from 'react'
import { Provider as ReduxProvider } from 'react-redux'
import { render } from 'react-snapshot'
import App from './App'
import createStore from './createStore'
import './index.css'
import * as serviceWorker from './serviceWorker'
import { theme } from '@chainlink/styleguide'

JavascriptTimeAgo.locale(en)
moment.defaultFormat = 'YYYY-MM-DD h:mm:ss A'

const muiTheme = createMuiTheme(theme)
const store = createStore()

render(
  <MuiThemeProvider theme={muiTheme}>
    <ReduxProvider store={store}>
      <App />
    </ReduxProvider>
  </MuiThemeProvider>,
  document.getElementById('root'),
)

serviceWorker.unregister()
