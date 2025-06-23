import {userReducer} from './userReducer'
import {
	authActionCreator,
	loginActionCreator,
	registerActionCreator,
	logoutActionCreator,
} from './userActionCreators'
import {
	selectUser,
	selectIsAuth,
	selectIsLoading,
	selectIsError,
	selectWillRedirect,
	selectBrowser,
	selectIsMobile,
	selectIsMiniapp,
	selectIsDesktop,
} from './userSelectors';
import {userSaga} from './userSagas';

export {
	logoutActionCreator,
	authActionCreator,
	userReducer,
	loginActionCreator,
	registerActionCreator,
	userSaga,
	selectUser,
	selectIsAuth,
	selectIsLoading,
	selectIsError,
	selectWillRedirect,
	selectBrowser,
	selectIsMobile,
	selectIsDesktop,
	selectIsMiniapp,
}