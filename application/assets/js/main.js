function test() {
    alert ("test")
}

// Travel Search Form JavaScript - Global Constants and Functions

// Element ID constants
const ID_DATABASE = 'database';
const ID_BOUNDS_INFO = 'bounds-info';
const ID_BOUNDS_X = 'bounds-x';
const ID_BOUNDS_Y = 'bounds-y';
const ID_TIME_BOUNDS_INFO = 'time-bounds-info';
const ID_BOUNDS_DEPARTURE = 'bounds-departure';
const ID_BOUNDS_ARRIVAL = 'bounds-arrival';
const ID_EXAMPLE_SOURCE = 'example-source';
const ID_EXAMPLE_DESTINATION = 'example-destination';
const ID_EXAMPLE_ARRIVAL_FROM = 'example-arrival-from';
const ID_EXAMPLE_ARRIVAL_TO = 'example-arrival-to';

/**
 * @param databaseName {string}
 * @returns {Promise<void>}
 */
async function fetchAndDisplayBounds(databaseName) {
    const boundsInfo = document.getElementById(ID_BOUNDS_INFO);
    const boundsX = document.getElementById(ID_BOUNDS_X);
    const boundsY = document.getElementById(ID_BOUNDS_Y);

    try {
        const response = await fetch(`/api/points/bounds?database=${encodeURIComponent(databaseName)}`);
        if (!response.ok) {
            throw new Error('Failed to fetch bounds');
        }
        const data = await response.json();

        if (data.bounds) {
            const bounds = data.bounds;
            boundsX.textContent = `${bounds.minX.toFixed(2)} to ${bounds.maxX.toFixed(2)}`;
            boundsY.textContent = `${bounds.minY.toFixed(2)} to ${bounds.maxY.toFixed(2)}`;

            // Show example coordinates in the correct syntax (x,y)
            document.getElementById(ID_EXAMPLE_SOURCE).textContent = `${bounds.minX.toFixed(2)},${bounds.minY.toFixed(2)}`;
            document.getElementById(ID_EXAMPLE_DESTINATION).textContent = `${bounds.maxX.toFixed(2)},${bounds.maxY.toFixed(2)}`;

            boundsInfo.style.display = 'block';
        } else {
            boundsInfo.style.display = 'none';
        }
    } catch (error) {
        console.error('Error fetching bounds:', error);
        boundsInfo.style.display = 'none';
    }
}

// Fetch and display travel time bounds when database changes
/**
 * @param databaseName {string}
 * @returns {Promise<void>}
 */
async function fetchAndDisplayTimeBounds(databaseName) {
    const timeBoundsInfo = document.getElementById(ID_TIME_BOUNDS_INFO);
    const boundsDeparture = document.getElementById(ID_BOUNDS_DEPARTURE);
    const boundsArrival = document.getElementById(ID_BOUNDS_ARRIVAL);

    try {
        const response = await fetch(`/api/travels/bounds?database=${encodeURIComponent(databaseName)}`);
        if (!response.ok) {
            throw new Error('Failed to fetch time bounds');
        }
        const data = await response.json();

        if (data.bounds) {
            const bounds = data.bounds;
            boundsDeparture.textContent = `${bounds.minDeparture} to ${bounds.maxDeparture}`;
            boundsArrival.textContent = `${bounds.minArrival} to ${bounds.maxArrival}`;

            // Show example dates in the correct format
            document.getElementById(ID_EXAMPLE_ARRIVAL_FROM).textContent = bounds.minArrival;
            document.getElementById(ID_EXAMPLE_ARRIVAL_TO).textContent = bounds.maxArrival;

            timeBoundsInfo.style.display = 'block';
        } else {
            timeBoundsInfo.style.display = 'none';
        }
    } catch (error) {
        console.error('Error fetching time bounds:', error);
        timeBoundsInfo.style.display = 'none';
    }
}

// Common date/time validation function
/**
 *
 * @param timeValueStr {string}
 * @returns {{valid: boolean, message: string}}
 */
function validateDateTimeFormat(timeValueStr) {
    if (!timeValueStr.trim()) {
        return { valid: false, message: 'This field is required' };
    }

    // Pattern for yyyy-mm-dd
    const datePattern = /^\d{4}-\d{2}-\d{2}$/;
    // Pattern for yyyy-mm-dd HH:ii
    const dateTimePattern = /^\d{4}-\d{2}-\d{2} \d{2}:\d{2}$/;

    if (datePattern.test(timeValueStr) || dateTimePattern.test(timeValueStr)) {
        // Additional validation: check if it's a valid date
        let dateStr = timeValueStr;
        if (datePattern.test(timeValueStr)) {
            dateStr = timeValueStr + ' 00:00';
        }

        const parts = dateStr.split(' ');
        const dateParts = parts[0].split('-');
        const timeParts = parts[1].split(':');

        const year = parseInt(dateParts[0]);
        const month = parseInt(dateParts[1]);
        const day = parseInt(dateParts[2]);
        const hour = parseInt(timeParts[0]);
        const minute = parseInt(timeParts[1]);

        // Basic range validation
        if (month < 1 || month > 12) {
            return { valid: false, message: 'Invalid month (must be 01-12)' };
        }
        if (day < 1 || day > 31) {
            return { valid: false, message: 'Invalid day (must be 01-31)' };
        }
        if (hour < 0 || hour > 23) {
            return { valid: false, message: 'Invalid hour (must be 00-23)' };
        }
        if (minute < 0 || minute > 59) {
            return { valid: false, message: 'Invalid minute (must be 00-59)' };
        }

        return { valid: true, message: '' };
    }

    return {
        valid: false,
        message: 'Invalid format. Use yyyy-mm-dd or yyyy-mm-dd HH:ii'
    };
}

// Setup validation for a date/time input field
/**
 *
 * @param inputElement {Element}
 * @param errorElement {Element}
 */
function setupDateTimeValidation(inputElement, errorElement) {
    inputElement.addEventListener('blur', () => {
        const result = validateDateTimeFormat(inputElement.value);
        if (!result.valid) {
            inputElement.classList.add('input-error');
            errorElement.textContent = result.message;
        } else {
            inputElement.classList.remove('input-error');
            errorElement.textContent = '';
        }
    });

    inputElement.addEventListener('input', () => {
        // Clear error while typing
        if (inputElement.classList.contains('input-error')) {
            inputElement.classList.remove('input-error');
            errorElement.textContent = '';
        }
    });
}

/**
 * Parse search value and determine parameter type
 *
 * @param searchValueStr {string}
 * @returns {{type: string, x: number, y: number}|{type: string, value: string}|null}
 */
function parseSearchValue(searchValueStr) {
    const trimmed = searchValueStr.trim();
    if (!trimmed) return null;

    const hasLetters = /[a-zA-Z]/.test(trimmed);
    const hasDigits = /[0-9]/.test(trimmed);
    const hasComma = /,/.test(trimmed);

    // Only digits and comma → X,Y coordinates
    if (!hasLetters && hasDigits && hasComma) {
        const parts = trimmed.split(',').map(p => p.trim());
        if (parts.length === 2) {
            const x = parseFloat(parts[0]);
            const y = parseFloat(parts[1]);
            if (!isNaN(x) && !isNaN(y)) {
                return { type: 'coordinates', x: x, y: y };
            }
        }
    }

    // Only letters → NamePart
    if (hasLetters && !hasDigits) {
        return { type: 'name', value: trimmed };
    }

    // Both letters and digits → IdPart
    if (hasLetters && hasDigits) {
        return { type: 'id', value: trimmed };
    }

    // Default to name search if unclear
    return { type: 'name', value: trimmed };
}

/**
 * Build API URL with query parameters
 * @param parsed {Object}
 * @returns {string}
 */
function buildApiUrl(parsed) {
    const baseUrl = '/api/points';
    const params = new URLSearchParams();
    params.append('limit', '20');

    if (parsed.type === 'coordinates') {
        params.append('x', parsed.x);
        params.append('y', parsed.y);
    } else if (parsed.type === 'name') {
        params.append('name_part', parsed.value);
    } else if (parsed.type === 'id') {
        params.append('id_part', parsed.value);
    }

    if ( parsed.database !== undefined) {
        params.append('database', parsed.database )
    }

    return `${baseUrl}?${params.toString()}`;
}

// Fetch points from API
async function fetchPoints(url) {
    try {
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error('Failed to fetch points');
        }
        const data = await response.json();
        return data.points || [];
    } catch (error) {
        console.error('Error fetching points:', error);
        return [];
    }
}

// Display points in source dropdown
function displayPoints(points) {
    const sourceDropdown = document.getElementById('source-dropdown');

    if (points.length === 0) {
        sourceDropdown.innerHTML = '<div class="loading">No points found</div>';
        sourceDropdown.classList.add('show');
        return;
    }

    sourceDropdown.innerHTML = '';
    points.forEach(point => {
        const item = document.createElement('div');
        item.className = 'dropdown-item';

        const nameDiv = document.createElement('div');
        nameDiv.className = 'dropdown-item-name';
        nameDiv.textContent = point.Name;

        const detailsDiv = document.createElement('div');
        detailsDiv.className = 'dropdown-item-details';
        detailsDiv.textContent = `ID: ${point.ID} | Coordinates: (${point.X.toFixed(2)}, ${point.Y.toFixed(2)})`;

        item.appendChild(nameDiv);
        item.appendChild(detailsDiv);

        item.dataset.id = point.ID;
        item.dataset.name = point.Name;
        item.dataset.x = point.X;
        item.dataset.y = point.Y;

        item.addEventListener('click', () => {
            selectSourcePoint(point);
        });

        sourceDropdown.appendChild(item);
    });
    sourceDropdown.classList.add('show');
}

// Select source point and update display
function selectSourcePoint(point) {
    const sourceHiddenInput = document.getElementById('source');
    const sourceSearchInput = document.getElementById('source-search');
    const sourceDropdown = document.getElementById('source-dropdown');

    sourceHiddenInput.value = point.ID;
    sourceSearchInput.value = '';
    sourceDropdown.classList.remove('show');

    document.getElementById('source-selected-name').textContent = point.Name;
    document.getElementById('source-selected-details').textContent =
        `ID: ${point.ID} | Coordinates: (${point.X.toFixed(2)}, ${point.Y.toFixed(2)})`;
    document.getElementById('source-selected').classList.add('show');
}

// Display points in destination dropdown
function displayDestinationPoints(points) {
    const destinationDropdown = document.getElementById('destination-dropdown');

    if (points.length === 0) {
        destinationDropdown.innerHTML = '<div class="loading">No points found</div>';
        destinationDropdown.classList.add('show');
        return;
    }

    destinationDropdown.innerHTML = '';
    points.forEach(point => {
        const item = document.createElement('div');
        item.className = 'dropdown-item';

        const nameDiv = document.createElement('div');
        nameDiv.className = 'dropdown-item-name';
        nameDiv.textContent = point.Name;

        const detailsDiv = document.createElement('div');
        detailsDiv.className = 'dropdown-item-details';
        detailsDiv.textContent = `ID: ${point.ID} | Coordinates: (${point.X.toFixed(2)}, ${point.Y.toFixed(2)})`;

        item.appendChild(nameDiv);
        item.appendChild(detailsDiv);

        item.dataset.id = point.ID;
        item.dataset.name = point.Name;
        item.dataset.x = point.X;
        item.dataset.y = point.Y;

        item.addEventListener('click', () => {
            selectDestinationPoint(point);
        });

        destinationDropdown.appendChild(item);
    });
    destinationDropdown.classList.add('show');
}

// Select destination point and update display
function selectDestinationPoint(point) {
    const destinationHiddenInput = document.getElementById('destination');
    const destinationSearchInput = document.getElementById('destination-search');
    const destinationDropdown = document.getElementById('destination-dropdown');

    destinationHiddenInput.value = point.ID;
    destinationSearchInput.value = '';
    destinationDropdown.classList.remove('show');

    document.getElementById('destination-selected-name').textContent = point.Name;
    document.getElementById('destination-selected-details').textContent =
        `ID: ${point.ID} | Coordinates: (${point.X.toFixed(2)}, ${point.Y.toFixed(2)})`;
    document.getElementById('destination-selected').classList.add('show');
}

// Initialize pre-selected points on page load
async function initializePreselectedPoints() {
    const databaseSelect = document.getElementById(ID_DATABASE);
    const sourceHiddenInput = document.getElementById('source');
    const destinationHiddenInput = document.getElementById('destination');
    const sourceId = sourceHiddenInput.value;
    const destinationId = destinationHiddenInput.value;

    if (!databaseSelect.value) {
        return;
    }

    // Load source point if pre-selected
    if (sourceId) {
        try {
            const params = new URLSearchParams();
            params.append('id_part', sourceId);
            params.append('limit', '1');
            params.append('database', databaseSelect.value);

            const response = await fetch(`/api/points?${params.toString()}`);
            if (response.ok) {
                const data = await response.json();
                if (data.points && data.points.length > 0) {
                    selectSourcePoint(data.points[0]);
                }
            }
        } catch (error) {
            console.error('Error loading source point:', error);
        }
    }

    // Load destination point if pre-selected
    if (destinationId) {
        try {
            const params = new URLSearchParams();
            params.append('id_part', destinationId);
            params.append('limit', '1');
            params.append('database', databaseSelect.value);

            const response = await fetch(`/api/points?${params.toString()}`);
            if (response.ok) {
                const data = await response.json();
                if (data.points && data.points.length > 0) {
                    selectDestinationPoint(data.points[0]);
                }
            }
        } catch (error) {
            console.error('Error loading destination point:', error);
        }
    }
}
