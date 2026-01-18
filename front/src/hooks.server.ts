import type { Handle } from '@sveltejs/kit';
import { paraglideMiddleware } from '$lib/paraglide/server';
import { sequence } from '@sveltejs/kit/hooks';

// Authentication middleware
const handleAuth: Handle = async ({ event, resolve }) => {
	// Try to get auth_token cookie
	const token = event.cookies.get('auth_token');

	if (token) {
		// Call backend /auth/me endpoint to validate token and get user
		try {
			const response = await fetch(`http://core:8080/auth/me`, {
				headers: {
					Cookie: `auth_token=${token}`
				}
			});

			if (response.ok) {
				event.locals.user = await response.json();
			} else {
				event.locals.user = null;
			}
		} catch (error) {
			console.error('Failed to validate auth token:', error);
			event.locals.user = null;
		}
	} else {
		event.locals.user = null;
	}

	// Protected routes check
	const protectedRoutes = ['/courses'];
	const isProtected = protectedRoutes.some((route) => event.url.pathname.startsWith(route));

	if (isProtected && !event.locals.user) {
		// Redirect to login page
		return new Response(null, {
			status: 303,
			headers: {
				location: '/login'
			}
		});
	}

	return resolve(event);
};

// Paraglide i18n middleware
const handleParaglide: Handle = ({ event, resolve }) =>
	paraglideMiddleware(event.request, ({ request, locale }) => {
		event.request = request;

		return resolve(event, {
			transformPageChunk: ({ html }) => html.replace('%paraglide.lang%', locale)
		});
	});

// Sequence the handles: auth first, then paraglide
export const handle: Handle = sequence(handleAuth, handleParaglide);
