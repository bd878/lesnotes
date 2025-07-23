import React from 'react';
import ReactDOM from 'react-dom/client';
import {mockTelegramEnv, emitEvent, isTMA, mockTelegramEnv, emitEvent, init, backButton, setDebug, bindThemeParamsCssVars, isThemeParamsCssVarsBound, themeParams} from '@telegram-apps/sdk-react';
import Tag from '../../components/Tag';
import Footer from '../../components/Footer';
import MiniappMainPage from '../../components/MiniappMainPage';
import StoreProvider from '../../providers/Store';
import i18n from '../../../i18n';

const noInsets = {
	left: 0,
	top: 0,
	bottom: 0,
	right: 0,
} as const;
const mockThemeParams = {
	accent_text_color: '#6ab2f2',
	bg_color: '#17212b',
	button_color: '#5288c1',
	button_text_color: '#ffffff',
	destructive_text_color: '#ec3942',
	header_bg_color: '#17212b',
	hint_color: '#708499',
	link_color: '#6ab3f3',
	secondary_bg_color: '#232e3c',
	section_bg_color: '#17212b',
	section_header_text_color: '#6ab3f3',
	subtitle_text_color: '#708499',
	text_color: '#f5f5f5',
} as const;

mockTelegramEnv({
	launchParams: {
		tgWebAppThemeParams: mockThemeParams,
		tgWebAppData: new URLSearchParams([
			['user', JSON.stringify({
				id: 1,
				first_name: 'Pavel',
			})],
			['hash', ''],
			['signature', ''],
			['auth_date', Date.now().toString()],
		]),
		tgWebAppStartParam: 'debug',
		tgWebAppVersion: '8',
		tgWebAppPlatform: 'tdesktop',
	},
	onEvent(e) {
		if (e[0] === 'web_app_request_theme') {
			return emitEvent('theme_changed', { theme_params: mockThemeParams });
		}
		if (e[0] === 'web_app_request_viewport') {
			return emitEvent('viewport_changed', {
				height: window.innerHeight,
				width: window.innerWidth,
				is_expanded: true,
				is_state_stable: true,
			});
		}
		if (e[0] === 'web_app_request_content_safe_area') {
			return emitEvent('content_safe_area_changed', noInsets);
		}
		if (e[0] === 'web_app_request_safe_area') {
			return emitEvent('safe_area_changed', noInsets);
		}
	},
})

setDebug(true);
init();

if (themeParams.mountSync.isAvailable()) {
	themeParams.mountSync()
	console.log("theme params mounted?:", themeParams.isMounted())
	bindThemeParamsCssVars()
	console.log("is theme params css bound?:", isThemeParamsCssVarsBound())
}
backButton.mount();

function Main() {
	if (!isTMA())
		return (
			<Tag css="h-screen flex flex-col items-start relative">
				{i18n("miniapp_only")}
				<Tag css="bg-white">
					<Footer />
				</Tag>
			</Tag>
		);

	return (
		<Tag css="h-screen flex flex-col items-start relative p-2 pb-4 bg-body">
			<StoreProvider isMiniapp={true}>
				<MiniappMainPage />
			</StoreProvider>

			<Footer textColor="text-subtitle" />
		</Tag>
	);
}

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Main />);
