function initLocationSelect1(provinces, districts, wards, $select_province, $select_district, $select_ward) {
    $select_province.append($('<option>', {
        value: 0,
        text: "--"
    }));
    $.each(provinces, function(index, province) {
        $select_province.append($('<option>', {
            value: province.code,
            text: province.name
        }));
    });

    $select_province.change(function() {
        var selectedProvinceCode = $(this).val();

        $select_ward.empty();

        $select_district.empty();
        $select_district.append($('<option>', {
            value: 0,
            text: "--"
        }));

        $.each(districts, function(index, district) {
            if(district.province_code === selectedProvinceCode) {
                $select_district.append($('<option>', {
                    value: district.code,
                    text: district.full_name
                }));
            }
        });
    });

    $select_district.change(function() {
        var selectedDistrictCode = $(this).val();
        $select_ward.empty();
        $select_ward.append($('<option>', {
            value: 0,
            text: "--"
        }));

        $.each(wards, function(index, ward) {
            if(ward.district_code === selectedDistrictCode) {
                $select_ward.append($('<option>', {
                    value: ward.code,
                    text: ward.full_name
                }));
            }
        });
    });
}

function initLocationSelect2(provinces, districts, wards, userProvinceCode, userDistrictCode, userWardCode, $select_province, $select_district, $select_ward) {
    $select_province.append($('<option>', {
        value: 0,
        text: "--"
    }));

    $.each(provinces, function(index, province) {
        const $option = $('<option>', {
            value: province.code,
            text: province.name
        });

        if (province.code === userProvinceCode) {
            $option.prop('selected', true);
        }

        $select_province.append($option);
    });

    renderDistricts(userProvinceCode);
    renderWards(userDistrictCode);

    function renderDistricts(provinceCode) {
        const selectedProvinceCode = provinceCode;
        $select_ward.empty();
        $select_district.empty();
        $select_district.append($('<option>', {
            value: 0,
            text: "--"
        }));

        $.each(districts, function(index, district) {
            if (district.province_code === selectedProvinceCode) {
                var $option = $('<option>', {
                    value: district.code,
                    text: district.full_name
                });
                if (district.code === userDistrictCode) {
                    $option.prop('selected', true);
                }
                $select_district.append($option);
            }
        });
    }

    function renderWards(districtCode) {
        var selectedDistrictCode = districtCode;
        $select_ward.empty();
        $select_ward.append($('<option>', {
            value: 0,
            text: "--"
        }));

        $.each(wards, function(index, ward) {
            if (ward.district_code === selectedDistrictCode) {
                var $option = $('<option>', {
                    value: ward.code,
                    text: ward.full_name
                });
                if (ward.code === userWardCode) {
                    $option.prop('selected', true);
                }
                $select_ward.append($option);
            }
        });
    }

    $select_province.change(function() {
        const selectedProvinceCode = $(this).val();
        if (selectedProvinceCode === userProvinceCode) {
            renderDistricts(selectedProvinceCode);
            renderWards(userDistrictCode);
        } else {
            renderDistricts(selectedProvinceCode);
        }
    });

    $select_district.change(function() {
        const selectedDistrictCode = $(this).val();
        renderWards(selectedDistrictCode);
    });
}