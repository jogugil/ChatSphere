export async function getClientInformation(): Promise<string> {
    // Obtener detalles del navegador
    const navegador = {
        userAgent: navigator.userAgent, // Cadena completa del navegador
        language: navigator.language,  // Idioma del navegador
        vendor: navigator.vendor,      // Proveedor del navegador
        urlActual: window.location.href, // URL actual del sitio
    };

    // Función para obtener la IP y país
    const getLocation = async (): Promise<{ ip: string; pais: string }> => {
        try {
            const respuesta = await fetch('https://ipapi.co/json/'); // Usa ipapi para obtener más detalles
            if (!respuesta.ok) {
                return { ip: 'Desconocida', pais: 'Desconocido' };
            }
            const datos = await respuesta.json();
            return { ip: datos.ip, pais: datos.country_name || 'Desconocido' }; // IP y país
        } catch (error: any) {
            // Evitar que se registre en la consola para el caso de bloqueadores
            if (error instanceof Error && error.message.includes("ERR_BLOCKED_BY_ADBLOCKER")) {
                console.log("Solicitud bloqueada por bloqueador de anuncios, se ignorará.");
            } else {
                // Solo loguear el error si es otro tipo de error
                console.log('Error al obtener la ubicación:', error.message || error);
            }
            return { ip: 'Desconocida', pais: 'Desconocido' };
        }
    };

    // Llamada para obtener IP y país
    const { ip, pais } = await getLocation();

    // Crear un objeto con toda la información del cliente
    const informacionCliente = {
        navegador: navegador,
        ip: ip,
        pais: pais,
    };

    // Convertir la información a JSON formateado
    return JSON.stringify(informacionCliente, null, 2);
}